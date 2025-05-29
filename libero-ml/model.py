"""
Poisson-based football score prediction model.
Predicts exact scores using two independent Poisson regression models for home and away goals.
Uses rolling team statistics and betting odds for enhanced accuracy.
"""

import pandas as pd
import numpy as np
from sklearn.linear_model import PoissonRegressor
from sklearn.preprocessing import LabelEncoder, StandardScaler
from sklearn.model_selection import TimeSeriesSplit
from sklearn.metrics import mean_absolute_error
import glob
import os
from datetime import datetime
import warnings
import random
import math
warnings.filterwarnings('ignore')

class SoccerPredictor:
    def __init__(self):
        # Ensemble of models for better predictions
        self.home_goals_model = PoissonRegressor(alpha=0.5, max_iter=1000)  # Reduced regularization
        self.away_goals_model = PoissonRegressor(alpha=0.5, max_iter=1000)
        
        # Add secondary models for ensemble
        self.home_goals_model_2 = PoissonRegressor(alpha=1.5, max_iter=1000)  # Higher regularization
        self.away_goals_model_2 = PoissonRegressor(alpha=1.5, max_iter=1000)
        
        # Feature scalers - separate for home and away
        self.home_scaler = StandardScaler()
        self.away_scaler = StandardScaler()
        
        # Label encoders for categorical variables
        self.league_encoder = LabelEncoder()
        self.team_encoder = LabelEncoder()
        
        # Team strength cache for dynamic calculations
        self.team_strengths = {}
        self.league_averages = {}
        
        # Data storage
        self.data = None
        self.home_features = None
        self.away_features = None
        
    def load_data(self, data_dir='data'):
        """Load and combine all historical match data from CSV files"""
        all_files = glob.glob(os.path.join(data_dir, '*.csv'))
        dataframes = []
        
        print(f"📁 Found {len(all_files)} data files")
        
        for file in all_files:
            try:
                df = pd.read_csv(file)
                # Convert Date to datetime
                df['Date'] = pd.to_datetime(df['Date'], format='%d/%m/%Y')
                dataframes.append(df)
                print(f"✅ Loaded {file}: {len(df)} matches")
            except Exception as e:
                print(f"❌ Error loading {file}: {e}")
                continue
                
        if not dataframes:
            raise ValueError("No valid data files found!")
            
        self.data = pd.concat(dataframes, ignore_index=True)
        # Sort by date to maintain temporal order for rolling calculations
        self.data = self.data.sort_values('Date').reset_index(drop=True)
        
        print(f"📊 Total dataset: {len(self.data)} matches")
        print(f"📅 Date range: {self.data['Date'].min()} to {self.data['Date'].max()}")
        print(f"🏆 Leagues: {', '.join(self.data['Div'].unique())}")
        
        return self.data
    
    def calculate_team_form(self, window=10):
        """Calculate comprehensive team statistics for the last N games"""
        print(f"📈 Calculating comprehensive team form (window={window} games)...")
        
        # Initialize columns for comprehensive rolling stats
        stats_cols = [
            # Overall team form (last 10 games regardless of venue)
            'home_team_overall_goals_avg', 'home_team_overall_conceded_avg', 'home_team_overall_form_points',
            'away_team_overall_goals_avg', 'away_team_overall_conceded_avg', 'away_team_overall_form_points',
            
            # Venue-specific form
            'home_team_home_goals_avg', 'home_team_home_conceded_avg', 'home_team_home_form_points',
            'away_team_away_goals_avg', 'away_team_away_conceded_avg', 'away_team_away_form_points',
            
            # Attack vs Defense ratings
            'home_team_attack_rating', 'home_team_defense_rating',
            'away_team_attack_rating', 'away_team_defense_rating',
            
            # Recent momentum (last 5 vs previous 5)
            'home_team_recent_momentum', 'away_team_recent_momentum',
            
            # Head-to-head record
            'h2h_home_advantage', 'h2h_goals_ratio'
        ]
        
        for col in stats_cols:
            self.data[col] = 0.0
        
        # Calculate for each match using only prior data
        for idx in range(len(self.data)):
            current_match = self.data.iloc[idx]
            current_date = current_match['Date']
            home_team = current_match['HomeTeam']
            away_team = current_match['AwayTeam']
            
            # Get all historical matches for both teams before current date
            home_team_matches = self.data[
                (self.data['Date'] < current_date) & 
                ((self.data['HomeTeam'] == home_team) | (self.data['AwayTeam'] == home_team))
            ].tail(window)
            
            away_team_matches = self.data[
                (self.data['Date'] < current_date) & 
                ((self.data['HomeTeam'] == away_team) | (self.data['AwayTeam'] == away_team))
            ].tail(window)
            
            # Get venue-specific matches
            home_team_home_matches = self.data[
                (self.data['Date'] < current_date) & 
                (self.data['HomeTeam'] == home_team)
            ].tail(window//2)  # Last 5 home games
            
            away_team_away_matches = self.data[
                (self.data['Date'] < current_date) & 
                (self.data['AwayTeam'] == away_team)
            ].tail(window//2)  # Last 5 away games
            
            # Calculate HOME TEAM overall form (last 10 games)
            if len(home_team_matches) > 0:
                home_goals, home_conceded, home_points = self._calculate_team_stats(home_team_matches, home_team)
                self.data.at[idx, 'home_team_overall_goals_avg'] = home_goals
                self.data.at[idx, 'home_team_overall_conceded_avg'] = home_conceded
                self.data.at[idx, 'home_team_overall_form_points'] = home_points
                
                # Attack and defense ratings (goals per game vs league average)
                league_avg_goals = 1.5  # Approximate league average
                self.data.at[idx, 'home_team_attack_rating'] = home_goals / league_avg_goals
                self.data.at[idx, 'home_team_defense_rating'] = league_avg_goals / max(home_conceded, 0.1)
            
            # Calculate AWAY TEAM overall form (last 10 games)
            if len(away_team_matches) > 0:
                away_goals, away_conceded, away_points = self._calculate_team_stats(away_team_matches, away_team)
                self.data.at[idx, 'away_team_overall_goals_avg'] = away_goals
                self.data.at[idx, 'away_team_overall_conceded_avg'] = away_conceded
                self.data.at[idx, 'away_team_overall_form_points'] = away_points
                
                # Attack and defense ratings
                self.data.at[idx, 'away_team_attack_rating'] = away_goals / league_avg_goals
                self.data.at[idx, 'away_team_defense_rating'] = league_avg_goals / max(away_conceded, 0.1)
            
            # Calculate HOME TEAM venue-specific form (at home)
            if len(home_team_home_matches) > 0:
                home_goals_h, home_conceded_h, home_points_h = self._calculate_venue_stats(
                    home_team_home_matches, home_team, 'home'
                )
                self.data.at[idx, 'home_team_home_goals_avg'] = home_goals_h
                self.data.at[idx, 'home_team_home_conceded_avg'] = home_conceded_h
                self.data.at[idx, 'home_team_home_form_points'] = home_points_h
            
            # Calculate AWAY TEAM venue-specific form (away)
            if len(away_team_away_matches) > 0:
                away_goals_a, away_conceded_a, away_points_a = self._calculate_venue_stats(
                    away_team_away_matches, away_team, 'away'
                )
                self.data.at[idx, 'away_team_away_goals_avg'] = away_goals_a
                self.data.at[idx, 'away_team_away_conceded_avg'] = away_conceded_a
                self.data.at[idx, 'away_team_away_form_points'] = away_points_a
            
            # Calculate recent momentum (last 5 vs previous 5)
            if len(home_team_matches) >= 10:
                recent_home = home_team_matches.tail(5)
                previous_home = home_team_matches.iloc[-10:-5]
                recent_points = self._calculate_team_stats(recent_home, home_team)[2]
                previous_points = self._calculate_team_stats(previous_home, home_team)[2]
                self.data.at[idx, 'home_team_recent_momentum'] = recent_points - previous_points
            
            if len(away_team_matches) >= 10:
                recent_away = away_team_matches.tail(5)
                previous_away = away_team_matches.iloc[-10:-5]
                recent_points = self._calculate_team_stats(recent_away, away_team)[2]
                previous_points = self._calculate_team_stats(previous_away, away_team)[2]
                self.data.at[idx, 'away_team_recent_momentum'] = recent_points - previous_points
            
            # Calculate head-to-head record
            h2h_matches = self.data[
                (self.data['Date'] < current_date) & 
                (((self.data['HomeTeam'] == home_team) & (self.data['AwayTeam'] == away_team)) |
                 ((self.data['HomeTeam'] == away_team) & (self.data['AwayTeam'] == home_team)))
            ].tail(5)  # Last 5 H2H matches
            
            if len(h2h_matches) > 0:
                h2h_home_wins = 0
                h2h_home_goals = 0
                h2h_away_goals = 0
                
                for _, h2h_match in h2h_matches.iterrows():
                    if h2h_match['HomeTeam'] == home_team:
                        # Home team was playing at home
                        if h2h_match['FTHG'] > h2h_match['FTAG']:
                            h2h_home_wins += 1
                        h2h_home_goals += h2h_match['FTHG']
                        h2h_away_goals += h2h_match['FTAG']
                    else:
                        # Home team was playing away
                        if h2h_match['FTAG'] > h2h_match['FTHG']:
                            h2h_home_wins += 1
                        h2h_home_goals += h2h_match['FTAG']
                        h2h_away_goals += h2h_match['FTHG']
                
                self.data.at[idx, 'h2h_home_advantage'] = h2h_home_wins / len(h2h_matches)
                self.data.at[idx, 'h2h_goals_ratio'] = h2h_home_goals / max(h2h_away_goals, 0.1)
                
        print(f"✅ Comprehensive team form calculation completed")
    
    def _calculate_team_stats(self, matches, team_name):
        """Calculate goals, conceded, and points for a team across matches"""
        if len(matches) == 0:
            return 0.0, 0.0, 0.0
            
        goals = 0
        conceded = 0
        points = 0
        
        for _, match in matches.iterrows():
            if match['HomeTeam'] == team_name:
                # Team played at home
                goals += match['FTHG']
                conceded += match['FTAG']
                if match['FTHG'] > match['FTAG']:
                    points += 3
                elif match['FTHG'] == match['FTAG']:
                    points += 1
            else:
                # Team played away
                goals += match['FTAG']
                conceded += match['FTHG']
                if match['FTAG'] > match['FTHG']:
                    points += 3
                elif match['FTAG'] == match['FTHG']:
                    points += 1
        
        return goals / len(matches), conceded / len(matches), points / len(matches)
    
    def _calculate_venue_stats(self, matches, team_name, venue):
        """Calculate venue-specific stats (home or away)"""
        if len(matches) == 0:
            return 0.0, 0.0, 0.0
            
        goals = 0
        conceded = 0
        points = 0
        
        for _, match in matches.iterrows():
            if venue == 'home' and match['HomeTeam'] == team_name:
                goals += match['FTHG']
                conceded += match['FTAG']
                if match['FTHG'] > match['FTAG']:
                    points += 3
                elif match['FTHG'] == match['FTAG']:
                    points += 1
            elif venue == 'away' and match['AwayTeam'] == team_name:
                goals += match['FTAG']
                conceded += match['FTHG']
                if match['FTAG'] > match['FTHG']:
                    points += 3
                elif match['FTAG'] == match['FTHG']:
                    points += 1
        
        return goals / len(matches), conceded / len(matches), points / len(matches)
    
    def engineer_features(self):
        """Create features for Poisson regression models"""
        print("🔧 Engineering features for Poisson models...")
        
        # Calculate rolling team form
        self.calculate_team_form()
        
        # Encode categorical variables
        self.data['HomeTeam_encoded'] = self.team_encoder.fit_transform(self.data['HomeTeam'])
        self.data['AwayTeam_encoded'] = self.team_encoder.transform(self.data['AwayTeam'])
        self.data['League_encoded'] = self.league_encoder.fit_transform(self.data['Div'])
        
        # Time-based features
        self.data['DayOfWeek'] = self.data['Date'].dt.dayofweek
        self.data['Month'] = self.data['Date'].dt.month
        self.data['IsWeekend'] = (self.data['Date'].dt.dayofweek >= 5).astype(int)
        
        # Enhanced betting odds features from multiple bookmakers
        betting_houses = ['B365', 'BW', 'IW', 'PS', 'WH', 'VC']
        valid_odds_cols = []
        
        for house in betting_houses:
            home_col = f'{house}H'
            draw_col = f'{house}D' 
            away_col = f'{house}A'
            
            if all(col in self.data.columns for col in [home_col, draw_col, away_col]):
                valid_odds_cols.extend([home_col, draw_col, away_col])
                
                # Convert odds to implied probabilities
                self.data[f'prob_{house.lower()}_home'] = 1 / self.data[home_col]
                self.data[f'prob_{house.lower()}_draw'] = 1 / self.data[draw_col]
                self.data[f'prob_{house.lower()}_away'] = 1 / self.data[away_col]
                
                # Normalize probabilities (remove bookmaker margin)
                total_prob = (self.data[f'prob_{house.lower()}_home'] + 
                             self.data[f'prob_{house.lower()}_draw'] + 
                             self.data[f'prob_{house.lower()}_away'])
                self.data[f'prob_{house.lower()}_home_norm'] = self.data[f'prob_{house.lower()}_home'] / total_prob
                self.data[f'prob_{house.lower()}_away_norm'] = self.data[f'prob_{house.lower()}_away'] / total_prob
        
        # Average betting market consensus
        if valid_odds_cols:
            home_prob_cols = [col for col in self.data.columns if 'prob_' in col and '_home_norm' in col]
            away_prob_cols = [col for col in self.data.columns if 'prob_' in col and '_away_norm' in col]
            
            if home_prob_cols and away_prob_cols:
                self.data['market_home_prob'] = self.data[home_prob_cols].mean(axis=1)
                self.data['market_away_prob'] = self.data[away_prob_cols].mean(axis=1)
                self.data['market_draw_prob'] = 1 - self.data['market_home_prob'] - self.data['market_away_prob']
                
                # Ensure valid probability ranges and handle NaN
                self.data['market_home_prob'] = self.data['market_home_prob'].fillna(0.4).clip(0.1, 0.8)
                self.data['market_away_prob'] = self.data['market_away_prob'].fillna(0.3).clip(0.1, 0.8) 
                self.data['market_draw_prob'] = self.data['market_draw_prob'].fillna(0.3).clip(0.1, 0.8)
        
        # Over/Under market intelligence 
        if 'B365>2.5' in self.data.columns and 'B365<2.5' in self.data.columns:
            self.data['market_over_2_5_prob'] = 1 / self.data['B365>2.5']
            self.data['market_under_2_5_prob'] = 1 / self.data['B365<2.5']
            # Normalize
            total_ou = self.data['market_over_2_5_prob'] + self.data['market_under_2_5_prob']
            self.data['market_over_2_5_norm'] = self.data['market_over_2_5_prob'] / total_ou
            
        # Shot-based features (if available)
        if 'HS' in self.data.columns and 'AS' in self.data.columns:
            self.data['shot_ratio'] = self.data['HS'] / (self.data['AS'] + 1)  # Avoid division by zero
            self.data['total_shots'] = self.data['HS'] + self.data['AS']
            
        if 'HST' in self.data.columns and 'AST' in self.data.columns:
            self.data['shots_on_target_ratio'] = self.data['HST'] / (self.data['AST'] + 1)
            self.data['total_shots_on_target'] = self.data['HST'] + self.data['AST']
            
        # Discipline features (fouls and cards)
        if 'HF' in self.data.columns and 'AF' in self.data.columns:
            self.data['foul_ratio'] = self.data['HF'] / (self.data['AF'] + 1)
            
        if 'HY' in self.data.columns and 'AY' in self.data.columns:
            self.data['yellow_card_ratio'] = self.data['HY'] / (self.data['AY'] + 1)
            
        # Corner kicks (indication of attacking pressure)
        if 'HC' in self.data.columns and 'AC' in self.data.columns:
            self.data['corner_ratio'] = self.data['HC'] / (self.data['AC'] + 1)
            self.data['total_corners'] = self.data['HC'] + self.data['AC']
        
        print("📊 Enhanced betting odds and match statistics features engineered")
        
        # Define feature sets for home and away goal prediction
        base_features = ['League_encoded', 'DayOfWeek', 'Month', 'IsWeekend']
        
        # Enhanced feature sets with new data
        market_features = []
        if 'market_home_prob' in self.data.columns:
            market_features.extend(['market_home_prob', 'market_away_prob', 'market_draw_prob'])
        if 'market_over_2_5_norm' in self.data.columns:
            market_features.append('market_over_2_5_norm')
            
        match_stats_features = []
        for col in ['shot_ratio', 'shots_on_target_ratio', 'foul_ratio', 'yellow_card_ratio', 'corner_ratio']:
            if col in self.data.columns:
                match_stats_features.append(col)
        
        # Features for predicting HOME goals
        self.home_features = base_features + [
            'HomeTeam_encoded', 'AwayTeam_encoded',
            'home_team_overall_goals_avg', 'home_team_overall_conceded_avg', 'home_team_overall_form_points',
            'home_team_home_goals_avg', 'home_team_home_conceded_avg', 'home_team_home_form_points',
            'home_team_attack_rating', 'home_team_defense_rating',
            'home_team_recent_momentum',
            'away_team_overall_conceded_avg', 'away_team_defense_rating',  # Opposition defense
            'h2h_home_advantage', 'h2h_goals_ratio'  # Head-to-head record
        ] + market_features + match_stats_features
        
        # Features for predicting AWAY goals  
        self.away_features = base_features + [
            'HomeTeam_encoded', 'AwayTeam_encoded',
            'away_team_overall_goals_avg', 'away_team_overall_conceded_avg', 'away_team_overall_form_points',
            'away_team_away_goals_avg', 'away_team_away_conceded_avg', 'away_team_away_form_points',
            'away_team_attack_rating', 'away_team_defense_rating',
            'away_team_recent_momentum',
            'home_team_overall_conceded_avg', 'home_team_defense_rating',  # Opposition defense
            'h2h_home_advantage', 'h2h_goals_ratio'  # Head-to-head record (inverse perspective)
        ] + market_features + match_stats_features
        
        # Filter features that actually exist in the data
        self.home_features = [f for f in self.home_features if f in self.data.columns]
        self.away_features = [f for f in self.away_features if f in self.data.columns]
        
        print(f"🏠 Home goal prediction features ({len(self.home_features)}): {len([f for f in self.home_features if 'market_' in f])} market + {len([f for f in self.home_features if any(x in f for x in ['shot', 'foul', 'card', 'corner'])])} match stats + {len(self.home_features) - len([f for f in self.home_features if 'market_' in f or any(x in f for x in ['shot', 'foul', 'card', 'corner'])])} team/form features")
        print(f"✈️  Away goal prediction features ({len(self.away_features)}): {len([f for f in self.away_features if 'market_' in f])} market + {len([f for f in self.away_features if any(x in f for x in ['shot', 'foul', 'card', 'corner'])])} match stats + {len(self.away_features) - len([f for f in self.away_features if 'market_' in f or any(x in f for x in ['shot', 'foul', 'card', 'corner'])])} team/form features")
        
        return self.home_features, self.away_features
    
    def prepare_training_data(self):
        """Prepare data for training with temporal validation"""
        self.engineer_features()
        
        # Remove early matches without sufficient form data
        valid_mask = (
            (self.data['home_team_overall_goals_avg'] > 0) & 
            (self.data['away_team_overall_goals_avg'] > 0)
        )
        valid_data = self.data[valid_mask].copy()
        
        print(f"📊 Valid matches for training: {len(valid_data)} (filtered from {len(self.data)})")
        
        # Prepare feature matrices with proper NaN handling
        X_home = valid_data[self.home_features].copy()
        X_away = valid_data[self.away_features].copy()
        
        # Fill NaN values with appropriate defaults
        print("🔧 Handling missing values...")
        
        # Fill categorical features with mode or 0
        categorical_cols = ['League_encoded', 'HomeTeam_encoded', 'AwayTeam_encoded']
        for col in categorical_cols:
            if col in X_home.columns:
                X_home[col] = X_home[col].fillna(0)
            if col in X_away.columns:
                X_away[col] = X_away[col].fillna(0)
        
        # Fill numerical features with median or sensible defaults
        for col in X_home.columns:
            if col not in categorical_cols:
                if 'prob' in col or 'market' in col:
                    # Betting odds - fill with neutral values
                    X_home[col] = X_home[col].fillna(0.33)
                elif 'ratio' in col:
                    # Ratios - fill with 1.0 (neutral)
                    X_home[col] = X_home[col].fillna(1.0)
                elif 'goals' in col or 'conceded' in col:
                    # Goals - fill with league average
                    X_home[col] = X_home[col].fillna(1.4)
                elif 'points' in col:
                    # Points - fill with average
                    X_home[col] = X_home[col].fillna(1.3)
                elif 'rating' in col:
                    # Ratings - fill with 1.0 (neutral)
                    X_home[col] = X_home[col].fillna(1.0)
                else:
                    # Other numerical - fill with median or 0
                    median_val = X_home[col].median()
                    X_home[col] = X_home[col].fillna(median_val if not pd.isna(median_val) else 0.0)
        
        for col in X_away.columns:
            if col not in categorical_cols:
                if 'prob' in col or 'market' in col:
                    X_away[col] = X_away[col].fillna(0.33)
                elif 'ratio' in col:
                    X_away[col] = X_away[col].fillna(1.0)
                elif 'goals' in col or 'conceded' in col:
                    X_away[col] = X_away[col].fillna(1.4)
                elif 'points' in col:
                    X_away[col] = X_away[col].fillna(1.3)
                elif 'rating' in col:
                    X_away[col] = X_away[col].fillna(1.0)
                else:
                    median_val = X_away[col].median()
                    X_away[col] = X_away[col].fillna(median_val if not pd.isna(median_val) else 0.0)
        
        # Final check for any remaining NaN values
        if X_home.isnull().any().any():
            print("⚠️  Warning: NaN values found in home features, filling with 0")
            X_home = X_home.fillna(0)
        
        if X_away.isnull().any().any():
            print("⚠️  Warning: NaN values found in away features, filling with 0")
            X_away = X_away.fillna(0)
        
        # Target variables (actual goals scored)
        y_home = valid_data['FTHG'].values  # Full-time home goals
        y_away = valid_data['FTAG'].values  # Full-time away goals
        
        # Scale numerical features (preserve categorical encodings)
        categorical_cols = ['League_encoded', 'HomeTeam_encoded', 'AwayTeam_encoded']
        
        # Scale home goal features
        home_numerical = [col for col in self.home_features if col not in categorical_cols]
        if home_numerical:
            print(f"🔢 Scaling {len(home_numerical)} numerical home features")
            X_home[home_numerical] = self.home_scaler.fit_transform(X_home[home_numerical])
        
        # Scale away goal features using away scaler
        away_numerical = [col for col in self.away_features if col not in categorical_cols]
        if away_numerical:
            print(f"🔢 Scaling {len(away_numerical)} numerical away features")
            X_away[away_numerical] = self.away_scaler.fit_transform(X_away[away_numerical])
        
        # Time-series split for proper temporal validation
        tscv = TimeSeriesSplit(n_splits=3)
        splits = list(tscv.split(X_home))
        train_idx, test_idx = splits[-1]  # Use the latest split
        
        return (X_home.iloc[train_idx], X_home.iloc[test_idx], 
                X_away.iloc[train_idx], X_away.iloc[test_idx],
                y_home[train_idx], y_home[test_idx],
                y_away[train_idx], y_away[test_idx])
    
    def train(self):
        """Train ensemble Poisson regression models with dynamic team strengths"""
        print("🚀 Training enhanced Poisson regression models...")
        
        # Calculate dynamic team strengths first
        self.calculate_dynamic_team_strengths()
        
        # Prepare training data
        (X_home_train, X_home_test, X_away_train, X_away_test,
         y_home_train, y_home_test, y_away_train, y_away_test) = self.prepare_training_data()
        
        print(f"📊 Training set size: {len(X_home_train)} matches")
        print(f"📊 Test set size: {len(X_home_test)} matches")
        
        # Train primary models (lower regularization for more variance)
        print("🏠 Training primary home goals Poisson model...")
        self.home_goals_model.fit(X_home_train, y_home_train)
        
        print("✈️  Training primary away goals Poisson model...")
        self.away_goals_model.fit(X_away_train, y_away_train)
        
        # Train secondary models (higher regularization for stability)
        print("🏠 Training secondary home goals Poisson model...")
        self.home_goals_model_2.fit(X_home_train, y_home_train)
        
        print("✈️  Training secondary away goals Poisson model...")
        self.away_goals_model_2.fit(X_away_train, y_away_train)
        
        # Evaluate models
        home_pred_train = self.home_goals_model.predict(X_home_train)
        home_pred_test = self.home_goals_model.predict(X_home_test)
        away_pred_train = self.away_goals_model.predict(X_away_train)
        away_pred_test = self.away_goals_model.predict(X_away_test)
        
        # Evaluate secondary models
        home_pred_train_2 = self.home_goals_model_2.predict(X_home_train)
        home_pred_test_2 = self.home_goals_model_2.predict(X_home_test)
        away_pred_train_2 = self.away_goals_model_2.predict(X_away_train)
        away_pred_test_2 = self.away_goals_model_2.predict(X_away_test)
        
        # Performance metrics
        print(f"\n📈 Primary Model Performance:")
        print(f"🏠 Home Goals - Train MAE: {mean_absolute_error(y_home_train, home_pred_train):.3f}")
        print(f"🏠 Home Goals - Test MAE: {mean_absolute_error(y_home_test, home_pred_test):.3f}")
        print(f"✈️  Away Goals - Train MAE: {mean_absolute_error(y_away_train, away_pred_train):.3f}")
        print(f"✈️  Away Goals - Test MAE: {mean_absolute_error(y_away_test, away_pred_test):.3f}")
        
        print(f"\n📈 Secondary Model Performance:")
        print(f"🏠 Home Goals - Train MAE: {mean_absolute_error(y_home_train, home_pred_train_2):.3f}")
        print(f"🏠 Home Goals - Test MAE: {mean_absolute_error(y_home_test, home_pred_test_2):.3f}")
        print(f"✈️  Away Goals - Train MAE: {mean_absolute_error(y_away_train, away_pred_train_2):.3f}")
        print(f"✈️  Away Goals - Test MAE: {mean_absolute_error(y_away_test, away_pred_test_2):.3f}")
        
        # Expected goals comparison
        print(f"\n📊 Expected Goals Comparison:")
        print(f"🏠 Home Goals - Train Mean: {y_home_train.mean():.2f} | Primary: {home_pred_train.mean():.2f} | Secondary: {home_pred_train_2.mean():.2f}")
        print(f"🏠 Home Goals - Test Mean: {y_home_test.mean():.2f} | Primary: {home_pred_test.mean():.2f} | Secondary: {home_pred_test_2.mean():.2f}")
        print(f"✈️  Away Goals - Train Mean: {y_away_train.mean():.2f} | Primary: {away_pred_train.mean():.2f} | Secondary: {away_pred_train_2.mean():.2f}")
        print(f"✈️  Away Goals - Test Mean: {y_away_test.mean():.2f} | Primary: {away_pred_test.mean():.2f} | Secondary: {away_pred_test_2.mean():.2f}")
        
        print("✅ Enhanced ensemble Poisson models training completed!")
    
    def predict_match(self, league, home_team, away_team, stats=None):
        """Enhanced prediction with improved team strength utilization and realistic score generation"""
        try:
            print(f"\n🔮 Enhanced Match Prediction: {home_team} vs {away_team}")
            
            # Get enhanced team strengths with all the new detailed stats
            home_strength = self.team_strengths.get(home_team, {
                'attack': 1.0, 'defense': 1.0, 'overall': 1.0, 'form': 1.0, 
                'confidence': 0.5, 'league': league, 'style': 'balanced', 
                'home_strength': 1.0, 'momentum': 0.0, 'win_rate': 0.33,
                'goals_per_game': 1.3, 'goals_conceded_per_game': 1.3,
                'home_goals_per_game': 1.4, 'quality_factor': 1.0
            })
            away_strength = self.team_strengths.get(away_team, {
                'attack': 1.0, 'defense': 1.0, 'overall': 1.0, 'form': 1.0, 
                'confidence': 0.5, 'league': league, 'style': 'balanced', 
                'home_strength': 1.0, 'momentum': 0.0, 'win_rate': 0.33,
                'goals_per_game': 1.3, 'goals_conceded_per_game': 1.3,
                'away_goals_per_game': 1.1, 'quality_factor': 1.0
            })
            
            print(f"🏠 {home_team} ({home_strength['style']})")
            print(f"   Attack: {home_strength['attack']:.2f} | Defense: {home_strength['defense']:.2f}")
            print(f"   Win Rate: {home_strength['win_rate']:.1%} | Goals/Game: {home_strength['goals_per_game']:.1f}")
            print(f"   Form: {home_strength['form']:.2f} | Momentum: {home_strength['momentum']:+.2f}")
            
            print(f"✈️  {away_team} ({away_strength['style']})")
            print(f"   Attack: {away_strength['attack']:.2f} | Defense: {away_strength['defense']:.2f}")
            print(f"   Win Rate: {away_strength['win_rate']:.1%} | Goals/Game: {away_strength['goals_per_game']:.1f}")
            print(f"   Form: {away_strength['form']:.2f} | Momentum: {away_strength['momentum']:+.2f}")
            
            # League context and cross-league adjustments
            home_league = home_strength.get('league', league)
            away_league = away_strength.get('league', league)
            
            # Enhanced cross-league handling with more sophisticated adjustments
            league_adjustment = 1.0
            if home_league != away_league:
                print(f"🌍 Cross-league match: {home_league} vs {away_league}")
                league_quality_factors = {
                    'E0': 1.00,   # Premier League
                    'SP1': 1.00,  # La Liga - same level as PL
                    'D1': 0.99,   # Bundesliga 
                    'I1': 0.99,   # Serie A - same level as Bundesliga
                    'F1': 0.98    # Ligue 1
                }
                home_factor = league_quality_factors.get(home_league, 1.0)
                away_factor = league_quality_factors.get(away_league, 1.0)
                league_adjustment = home_factor / away_factor
                league_adjustment = max(0.98, min(1.02, league_adjustment))  # Much smaller differences now
            
            # Get enhanced league context
            league_avg = self.league_averages.get(league, {
                'home_goals': 1.55, 'away_goals': 1.25, 'total_goals': 2.8, 
                'quality_factor': 1.0, 'competitiveness': 1.0
            })
            
            # More realistic strength-based calculation using actual team data
            home_attack_rating = home_strength['attack']
            away_attack_rating = away_strength['attack'] 
            home_defense_rating = home_strength['defense']
            away_defense_rating = away_strength['defense']
            
            # Quality factor adjustments - good teams should be more predictable
            home_quality = home_strength['quality_factor']
            away_quality = away_strength['quality_factor']
            
            # Tactical matchup analysis with more sophisticated interactions
            tactical_factor_home = 1.0
            tactical_factor_away = 1.0
            
            if home_strength['style'] == 'attacking' and away_strength['style'] == 'defensive':
                # Attacking team vs defensive - often neutralizes
                tactical_factor_home = 0.85
                tactical_factor_away = 1.10
                print("⚔️  Tactical: Attacking vs Defensive - neutralizing effect")
            elif home_strength['style'] == 'defensive' and away_strength['style'] == 'attacking':
                # Defensive vs attacking - home advantage helps defensive team
                tactical_factor_home = 1.10
                tactical_factor_away = 0.90
                print("⚔️  Tactical: Defensive vs Attacking - home defense advantage")
            elif home_strength['style'] == 'attacking' and away_strength['style'] == 'attacking':
                # Both attacking - high-scoring potential
                tactical_factor_home = 1.15
                tactical_factor_away = 1.15
                print("⚔️  Tactical: Both Attacking - high-scoring expected")
            elif home_strength['style'] == 'complete' or away_strength['style'] == 'complete':
                # Complete teams are more adaptable
                if home_strength['style'] == 'complete':
                    tactical_factor_home = 1.08
                if away_strength['style'] == 'complete':
                    tactical_factor_away = 1.05
                print("⚔️  Tactical: Complete team(s) - tactical advantage")
            
            # Enhanced home advantage calculation
            base_home_advantage = 1.12  # Base home advantage
            venue_strength_bonus = 0.03 * (home_strength['home_strength'] - 1.0)  # Venue-specific bonus
            home_advantage_factor = base_home_advantage + venue_strength_bonus
            away_adjustment = 0.88 - 0.02 * (away_strength['home_strength'] - 1.0)  # Teams good at home travel better
            
            # Form and momentum impact - more significant for realistic differences
            home_form_impact = 0.85 + 0.30 * (home_strength['form'] - 1.0)
            away_form_impact = 0.85 + 0.30 * (away_strength['form'] - 1.0)
            
            # Momentum creates swing factor
            momentum_home = 1.0 + home_strength['momentum'] * 0.25
            momentum_away = 1.0 + away_strength['momentum'] * 0.25
            
            # Calculate expected goals using enhanced team data and realistic interactions
            base_home_goals = (
                league_avg['home_goals'] * 
                home_attack_rating * 
                (2.4 - away_defense_rating * 0.8) *  # Defense impact (stronger effect)
                home_advantage_factor * 
                home_form_impact * 
                momentum_home * 
                tactical_factor_home *
                league_adjustment * 
                league_avg['quality_factor'] *
                home_quality
            )
            
            base_away_goals = (
                league_avg['away_goals'] * 
                away_attack_rating * 
                (2.4 - home_defense_rating * 0.8) * 
                away_adjustment * 
                away_form_impact * 
                momentum_away * 
                tactical_factor_away *
                (1.0 / league_adjustment) * 
                league_avg['quality_factor'] *
                away_quality
            )
            
            # Apply realistic bounds with wider range for variety
            base_home_goals = max(0.3, min(5.0, base_home_goals))
            base_away_goals = max(0.3, min(5.0, base_away_goals))
            
            print(f"📊 Strength-based prediction: {base_home_goals:.2f} - {base_away_goals:.2f}")
            
            # Model predictions with enhanced features (if models are trained)
            model_home_goals = base_home_goals
            model_away_goals = base_away_goals
            
            try:
                league_encoded = 0
                try:
                    league_encoded = self.league_encoder.transform([league])[0]
                except:
                    pass
                    
                home_encoded = self._find_team_encoding(home_team)
                away_encoded = self._find_team_encoding(away_team)
                
                current_date = datetime.now()
                
                # Enhanced feature engineering using all our team strength data
                enhanced_features = {
                    'League_encoded': league_encoded,
                    'HomeTeam_encoded': home_encoded,
                    'AwayTeam_encoded': away_encoded,
                    'DayOfWeek': current_date.weekday(),
                    'Month': current_date.month,
                    'IsWeekend': 1 if current_date.weekday() >= 5 else 0,
                    
                    # Enhanced team metrics from our detailed analysis
                    'home_team_overall_goals_avg': home_strength['goals_per_game'],
                    'home_team_overall_conceded_avg': home_strength['goals_conceded_per_game'],
                    'home_team_attack_rating': home_strength['attack'],
                    'home_team_defense_rating': home_strength['defense'],
                    'home_team_overall_form_points': home_strength['win_rate'] * 3,
                    'home_team_recent_momentum': home_strength['momentum'],
                    'home_team_home_goals_avg': home_strength.get('home_goals_per_game', home_strength['goals_per_game']),
                    
                    'away_team_overall_goals_avg': away_strength['goals_per_game'],
                    'away_team_overall_conceded_avg': away_strength['goals_conceded_per_game'],
                    'away_team_attack_rating': away_strength['attack'],
                    'away_team_defense_rating': away_strength['defense'],
                    'away_team_overall_form_points': away_strength['win_rate'] * 3,
                    'away_team_recent_momentum': away_strength['momentum'],
                    'away_team_away_goals_avg': away_strength.get('away_goals_per_game', away_strength['goals_per_game']),
                    
                    # Tactical and contextual features
                    'h2h_home_advantage': 0.5 + 0.25 * (home_strength['overall'] - away_strength['overall']),
                    'h2h_goals_ratio': max(0.5, min(2.0, home_strength['attack'] / max(away_strength['defense'], 0.5))),
                    
                    # Realistic market probabilities based on team strengths
                    'market_home_prob': max(0.20, min(0.70, 0.40 + 0.20 * (home_strength['overall'] - away_strength['overall']))),
                    'market_away_prob': max(0.15, min(0.60, 0.25 + 0.20 * (away_strength['overall'] - home_strength['overall']))),
                }
                
                # Ensure market probabilities sum correctly
                home_prob = enhanced_features['market_home_prob']
                away_prob = enhanced_features['market_away_prob']
                if home_prob + away_prob > 0.85:
                    scale = 0.85 / (home_prob + away_prob)
                    enhanced_features['market_home_prob'] = home_prob * scale
                    enhanced_features['market_away_prob'] = away_prob * scale
                
                enhanced_features['market_draw_prob'] = 1.0 - enhanced_features['market_home_prob'] - enhanced_features['market_away_prob']
                enhanced_features['market_over_2_5_norm'] = max(0.30, min(0.80, 0.45 + 0.25 * (home_strength['attack'] + away_strength['attack'] - 2.0)))
                
                # Process through models if available
                if hasattr(self, 'home_features') and hasattr(self, 'away_features'):
                    X_home = np.array([[enhanced_features.get(f, 0) for f in self.home_features]])
                    X_away = np.array([[enhanced_features.get(f, 0) for f in self.away_features]])
                    
                    X_home = np.nan_to_num(X_home, nan=0.0)
                    X_away = np.nan_to_num(X_away, nan=0.0)
                    
                    # Scale features if scalers are available
                    categorical_cols = ['League_encoded', 'HomeTeam_encoded', 'AwayTeam_encoded']
                    home_numerical_idx = [i for i, f in enumerate(self.home_features) if f not in categorical_cols]
                    away_numerical_idx = [i for i, f in enumerate(self.away_features) if f not in categorical_cols]
                    
                    if home_numerical_idx and hasattr(self.home_scaler, 'mean_'):
                        try:
                            X_home[:, home_numerical_idx] = self.home_scaler.transform(X_home[:, home_numerical_idx])
                        except:
                            pass
                            
                    if away_numerical_idx and hasattr(self.away_scaler, 'mean_'):
                        try:
                            X_away[:, away_numerical_idx] = self.away_scaler.transform(X_away[:, away_numerical_idx])
                        except:
                            pass
                    
                    # Ensemble predictions with bounds
                    model_home_1 = max(0.2, min(6.0, self.home_goals_model.predict(X_home)[0]))
                    model_away_1 = max(0.2, min(6.0, self.away_goals_model.predict(X_away)[0]))
                    model_home_2 = max(0.2, min(6.0, self.home_goals_model_2.predict(X_home)[0]))
                    model_away_2 = max(0.2, min(6.0, self.away_goals_model_2.predict(X_away)[0]))
                    
                    model_home_goals = 0.6 * model_home_1 + 0.4 * model_home_2
                    model_away_goals = 0.6 * model_away_1 + 0.4 * model_away_2
                    
                    print(f"📈 Model ensemble: {model_home_goals:.2f} - {model_away_goals:.2f}")
                
            except Exception as e:
                print(f"⚠️  Model prediction failed: {e}, using strength-based")
                model_home_goals = base_home_goals
                model_away_goals = base_away_goals
            
            # Intelligent hybrid weighting based on confidence and team quality
            avg_confidence = (home_strength['confidence'] + away_strength['confidence']) / 2
            strength_weight = max(0.35, min(0.75, avg_confidence * 0.8))  # Higher confidence = more strength weight
            model_weight = 1.0 - strength_weight
            
            final_home_goals = strength_weight * base_home_goals + model_weight * model_home_goals
            final_away_goals = strength_weight * base_away_goals + model_weight * model_away_goals
            
            print(f"🎯 Hybrid prediction: {final_home_goals:.2f} - {final_away_goals:.2f}")
            print(f"⚖️  Weights - Strength: {strength_weight:.1%}, Model: {model_weight:.1%}")
            
            # Adaptive variance based on team characteristics and match context
            strength_gap = abs(home_strength['overall'] - away_strength['overall'])
            confidence_gap = abs(home_strength['confidence'] - away_strength['confidence'])
            
            # Base variance - teams with similar strengths are less predictable
            base_variance = 0.08 if strength_gap > 0.5 else 0.15
            
            # Style-based variance adjustments
            if home_strength['style'] in ['attacking', 'complete'] or away_strength['style'] in ['attacking', 'complete']:
                base_variance += 0.05  # More unpredictable
            if 'defensive' in [home_strength['style'], away_strength['style']]:
                base_variance -= 0.03  # More predictable (lower scoring)
            
            # Confidence-based variance
            variance_factor = max(0.03, min(0.20, base_variance * (1 - avg_confidence)))
            
            # Apply controlled randomness for realism
            home_variance = 1 + random.uniform(-variance_factor, variance_factor)
            away_variance = 1 + random.uniform(-variance_factor, variance_factor)
            
            expected_home_goals = final_home_goals * home_variance
            expected_away_goals = final_away_goals * away_variance
            
            # Final realistic bounds
            expected_home_goals = max(0.2, min(6.0, expected_home_goals))
            expected_away_goals = max(0.2, min(6.0, expected_away_goals))
            
            print(f"🎲 Final expected: {expected_home_goals:.2f} - {expected_away_goals:.2f}")
            
            # Enhanced score calculation with realistic distribution
            try:
                from scipy.stats import poisson
                
                score_probabilities = {}
                max_goals = 6  # Allow up to 6 goals for rare high-scoring games
                
                for h in range(max_goals + 1):
                    for a in range(max_goals + 1):
                        prob = poisson.pmf(h, expected_home_goals) * poisson.pmf(a, expected_away_goals)
                        score_probabilities[(h, a)] = prob
                
                # Find most likely score
                most_likely_score = max(score_probabilities.items(), key=lambda x: x[1])
                most_likely_home, most_likely_away = most_likely_score[0]
                
                # Calculate outcome probabilities with better precision
                home_win_prob = sum(prob for (h, a), prob in score_probabilities.items() if h > a)
                draw_prob = sum(prob for (h, a), prob in score_probabilities.items() if h == a)
                away_win_prob = sum(prob for (h, a), prob in score_probabilities.items() if h < a)
                
                # Normalize probabilities
                total_prob = home_win_prob + draw_prob + away_win_prob
                if total_prob > 0:
                    home_win_prob /= total_prob
                    draw_prob /= total_prob
                    away_win_prob /= total_prob
                
            except ImportError:
                # Enhanced fallback calculation using team strengths
                most_likely_home = max(0, min(6, int(round(expected_home_goals))))
                most_likely_away = max(0, min(6, int(round(expected_away_goals))))
                
                # More sophisticated probability calculation
                goal_diff = expected_home_goals - expected_away_goals
                strength_diff = home_strength['overall'] - away_strength['overall']
                
                # Base probabilities considering team strengths
                if goal_diff > 0.8 or strength_diff > 0.3:
                    home_win_prob, draw_prob, away_win_prob = 0.58, 0.22, 0.20
                elif goal_diff < -0.8 or strength_diff < -0.3:
                    home_win_prob, draw_prob, away_win_prob = 0.20, 0.22, 0.58
                elif abs(goal_diff) < 0.3 and abs(strength_diff) < 0.2:
                    home_win_prob, draw_prob, away_win_prob = 0.35, 0.30, 0.35  # Very even match
                else:
                    home_win_prob, draw_prob, away_win_prob = 0.42, 0.26, 0.32  # Slight home advantage
            
            # Determine final prediction
            if home_win_prob > max(draw_prob, away_win_prob):
                prediction = 1
            elif away_win_prob > draw_prob:
                prediction = -1
            else:
                prediction = 0
            
            print(f"🏆 Predicted Score: {most_likely_home}-{most_likely_away}")
            print(f"📊 Win Probabilities: Home {home_win_prob:.1%} | Draw {draw_prob:.1%} | Away {away_win_prob:.1%}")
            
            return {
                'prediction': prediction,
                'probabilities': {
                    'home_win': float(home_win_prob),
                    'draw': float(draw_prob),
                    'away_win': float(away_win_prob)
                },
                'expected_home_goals': float(expected_home_goals),
                'expected_away_goals': float(expected_away_goals),
                'most_likely_home_score': most_likely_home,
                'most_likely_away_score': most_likely_away
            }
            
        except Exception as e:
            print(f"❌ Enhanced prediction error: {e}")
            import traceback
            traceback.print_exc()
            
            # Realistic fallback with team considerations
            fallback_home = random.randint(0, 3)
            fallback_away = random.randint(0, 3)
            
            return {
                'prediction': random.choice([-1, 0, 1]),
                'probabilities': {'home_win': 0.40, 'draw': 0.28, 'away_win': 0.32},
                'expected_home_goals': 1.3 + random.uniform(-0.5, 0.8),
                'expected_away_goals': 1.1 + random.uniform(-0.5, 0.8),
                'most_likely_home_score': fallback_home,
                'most_likely_away_score': fallback_away
            }
    
    def _find_team_encoding(self, team_name):
        """Helper method to find team encoding with fallbacks"""
        # First try exact match
        try:
            return self.team_encoder.transform([team_name])[0]
        except:
            pass
        
        # Try common variations
        variations = [
            team_name.replace(' ', ''),
            team_name.replace('-', ' '), 
            team_name.replace('Borussia Dortmund', 'Dortmund'),
            team_name.replace('Bayern Munich', 'Bayern'),
        ]
        
        for variation in variations:
            try:
                return self.team_encoder.transform([variation])[0]
            except:
                continue
        
        # Partial matching
        all_teams = list(self.team_encoder.classes_)
        for i, existing_team in enumerate(all_teams):
            if team_name.lower() in existing_team.lower() or existing_team.lower() in team_name.lower():
                print(f"🔍 Found similar team: '{team_name}' -> '{existing_team}'")
                return i
        
        # Random fallback for diversity
        fallback = random.randint(0, len(all_teams) - 1)
        print(f"⚠️  Unknown team '{team_name}', using random encoding: {fallback}")
        return fallback
    
    def get_available_teams(self):
        """Get list of all teams in the dataset"""
        if self.data is None:
            return []
        
        home_teams = set(self.data['HomeTeam'].unique())
        away_teams = set(self.data['AwayTeam'].unique())
        all_teams = sorted(list(home_teams.union(away_teams)))
        
        return all_teams
    
    def get_available_leagues(self):
        """Get list of all leagues in the dataset"""
        if self.data is None:
            return []
        
        return sorted(self.data['Div'].unique().tolist())

    def calculate_dynamic_team_strengths(self):
        """Calculate meaningful team strengths with better differentiation and realistic predictions"""
        print("💪 Calculating enhanced team strengths with improved W/L/D analysis...")
        
        # Calculate league-specific baselines with more detailed analysis
        league_stats = {}
        for league in self.data['Div'].unique():
            league_data = self.data[self.data['Div'] == league]
            league_stats[league] = {
                'avg_home_goals': league_data['FTHG'].mean(),
                'avg_away_goals': league_data['FTAG'].mean(),
                'avg_total_goals': (league_data['FTHG'] + league_data['FTAG']).mean(),
                'home_win_rate': (league_data['FTHG'] > league_data['FTAG']).mean(),
                'draw_rate': (league_data['FTHG'] == league_data['FTAG']).mean(),
                'away_win_rate': (league_data['FTHG'] < league_data['FTAG']).mean(),
                'high_scoring_rate': ((league_data['FTHG'] + league_data['FTAG']) > 2.5).mean(),
                'clean_sheet_rate': ((league_data['FTHG'] == 0) | (league_data['FTAG'] == 0)).mean()
            }
        
        # Calculate team performance with much better differentiation
        all_teams = set(self.data['HomeTeam'].unique()) | set(self.data['AwayTeam'].unique())
        
        for team in all_teams:
            # Find team's league and recent performance (increased sample for stability)
            team_league = None
            team_matches = self.data[
                (self.data['HomeTeam'] == team) | (self.data['AwayTeam'] == team)
            ]
            
            if len(team_matches) > 0:
                # Get most common league for this team
                league_counts = team_matches['Div'].value_counts()
                team_league = league_counts.index[0]
            
            if not team_league or len(team_matches) < 8:  # Increased minimum games
                # Default for insufficient data teams with more conservative ranges
                self.team_strengths[team] = {
                    'attack': 1.0, 'defense': 1.0, 'overall': 1.0, 'form': 1.0, 
                    'confidence': 0.4, 'league': team_league or 'E0',
                    'style': 'balanced', 'home_strength': 1.0, 'momentum': 0.0,
                    'win_rate': 0.33, 'draw_rate': 0.27, 'goals_per_game': 1.3,
                    'goals_conceded_per_game': 1.3, 'clean_sheet_rate': 0.25
                }
                continue
            
            # Get more extensive recent form (last 30 games for better assessment)
            recent_matches = team_matches.tail(30)
            league_baseline = league_stats[team_league]
            
            # Calculate comprehensive performance metrics
            total_goals_for = 0
            total_goals_against = 0
            total_points = 0
            wins = 0
            draws = 0
            losses = 0
            home_goals_for = 0
            home_goals_against = 0
            away_goals_for = 0
            away_goals_against = 0
            clean_sheets = 0
            big_wins = 0  # Wins by 2+ goals
            big_losses = 0  # Losses by 2+ goals
            comeback_wins = 0  # Wins after being behind
            
            home_games = 0
            away_games = 0
            home_wins = 0
            away_wins = 0
            home_points = 0
            away_points = 0
            
            for _, match in recent_matches.iterrows():
                if match['HomeTeam'] == team:
                    # Playing at home
                    gf, ga = match['FTHG'], match['FTAG']
                    home_goals_for += gf
                    home_goals_against += ga
                    home_games += 1
                    
                    if gf > ga:
                        wins += 1
                        home_wins += 1
                        home_points += 3
                        if gf - ga >= 2:
                            big_wins += 1
                    elif gf == ga:
                        draws += 1
                        home_points += 1
                    else:
                        losses += 1
                        if ga - gf >= 2:
                            big_losses += 1
                    
                    if ga == 0:
                        clean_sheets += 1
                else:
                    # Playing away
                    gf, ga = match['FTAG'], match['FTHG']
                    away_goals_for += gf
                    away_goals_against += ga
                    away_games += 1
                    
                    if gf > ga:
                        wins += 1
                        away_wins += 1
                        away_points += 3
                        if gf - ga >= 2:
                            big_wins += 1
                    elif gf == ga:
                        draws += 1
                        away_points += 1
                    else:
                        losses += 1
                        if ga - gf >= 2:
                            big_losses += 1
                    
                    if ga == 0:
                        clean_sheets += 1
                
                total_goals_for += gf
                total_goals_against += ga
                
                if gf > ga:
                    total_points += 3
                elif gf == ga:
                    total_points += 1
            
            # Calculate detailed rates and performance metrics
            games_played = len(recent_matches)
            avg_goals_for = total_goals_for / games_played
            avg_goals_against = total_goals_against / games_played
            avg_points = total_points / games_played
            win_rate = wins / games_played
            draw_rate = draws / games_played
            loss_rate = losses / games_played
            clean_sheet_rate = clean_sheets / games_played
            
            # Home/away performance analysis
            home_goals_per_game = home_goals_for / max(home_games, 1) if home_games > 0 else avg_goals_for
            away_goals_per_game = away_goals_for / max(away_games, 1) if away_games > 0 else avg_goals_for
            home_conceded_per_game = home_goals_against / max(home_games, 1) if home_games > 0 else avg_goals_against
            away_conceded_per_game = away_goals_against / max(away_games, 1) if away_games > 0 else avg_goals_against
            
            home_win_rate = home_wins / max(home_games, 1) if home_games > 0 else win_rate
            away_win_rate = away_wins / max(away_games, 1) if away_games > 0 else win_rate
            
            # Calculate enhanced strength metrics with much wider meaningful ranges
            expected_goals_for = (league_baseline['avg_home_goals'] + league_baseline['avg_away_goals']) / 2
            expected_points = 1.5  # Average points per game
            expected_win_rate = 0.35  # League average
            
            # Attack strength: wider range 0.4-2.0 for much better differentiation
            attack_base = avg_goals_for / max(expected_goals_for, 0.8)
            attack_strength = 0.4 + 1.6 * min(max(attack_base, 0.2), 1.8)  # Range: 0.4-2.0
            
            # Defense strength: wider range 0.4-2.0 (higher = better defense = fewer goals conceded)
            defense_base = expected_goals_for / max(avg_goals_against, 0.3)
            defense_strength = 0.4 + 1.6 * min(max(defense_base, 0.2), 1.8)  # Range: 0.4-2.0
            
            # Form strength based on points per game: range 0.5-1.8
            form_base = avg_points / expected_points
            form_strength = 0.5 + 1.3 * min(max(form_base, 0.2), 1.5)  # Range: 0.5-1.8
            
            # Win rate factor for overall strength assessment
            win_rate_factor = 1.0 + (win_rate - expected_win_rate) * 2.0  # More aggressive scaling
            
            # Quality factor based on goal difference and consistency
            goal_diff_per_game = (total_goals_for - total_goals_against) / games_played
            consistency_factor = 1.0 - (big_losses / max(games_played, 1)) + (big_wins / max(games_played, 1))
            quality_factor = 1.0 + (goal_diff_per_game * 0.3) + (consistency_factor - 1.0) * 0.2
            
            # Apply league quality with more significant differences
            league_quality_multipliers = {
                'E0': 1.00,   # Premier League
                'SP1': 1.00,  # La Liga - same level as PL
                'D1': 0.99,   # Bundesliga 
                'I1': 0.99,   # Serie A - same level as Bundesliga
                'F1': 0.98    # Ligue 1
            }
            league_multiplier = league_quality_multipliers.get(team_league, 1.0)
            
            # Apply all adjustments with more aggressive scaling
            attack_strength *= (0.7 + 0.3 * league_multiplier) * quality_factor * (0.8 + 0.4 * win_rate_factor)
            defense_strength *= (0.7 + 0.3 * league_multiplier) * quality_factor * (0.8 + 0.4 * (1 + clean_sheet_rate))
            form_strength *= (0.7 + 0.3 * league_multiplier) * win_rate_factor
            
            # Determine playing style with more nuanced classification
            goals_ratio = avg_goals_for / max(avg_goals_against, 0.5)
            style = 'balanced'
            if avg_goals_for > expected_goals_for * 1.3:
                if avg_goals_against < expected_goals_for * 0.9:
                    style = 'complete'  # High scoring + good defense
                else:
                    style = 'attacking'  # High scoring but leaky defense
            elif avg_goals_against < expected_goals_for * 0.7:
                style = 'defensive'  # Low scoring but solid defense
            elif win_rate > 0.55:
                style = 'efficient'  # Good results without extreme stats
            
            # Recent momentum calculation (last 8 games vs previous 8)
            momentum = 0.0
            if games_played >= 16:
                recent_8 = recent_matches.tail(8)
                previous_8 = recent_matches.iloc[-16:-8]
                
                recent_points = self._calculate_team_stats(recent_8, team)[2]
                previous_points = self._calculate_team_stats(previous_8, team)[2]
                momentum = (recent_points - previous_points) / 3.0  # Scale to -1.0 to +1.0 range
            
            # Overall strength calculation with better weighting
            overall_strength = (
                attack_strength * 0.35 + 
                defense_strength * 0.35 + 
                form_strength * 0.20 + 
                (1.0 + momentum * 0.5) * 0.10  # Momentum influence
            )
            
            # Home venue strength (how much better at home vs away)
            home_advantage = 1.0
            if home_games > 3 and away_games > 3:
                home_performance = (home_goals_per_game - home_conceded_per_game + home_win_rate * 3)
                away_performance = (away_goals_per_game - away_conceded_per_game + away_win_rate * 3)
                home_advantage = max(0.8, min(1.4, home_performance / max(away_performance, 0.5)))
            
            # Confidence based on sample size, consistency, and league quality
            sample_confidence = min(1.0, games_played / 25.0)  # Full confidence at 25+ games
            consistency_confidence = 1.0 - min(0.4, abs(goal_diff_per_game) / 3.0)  # Lower confidence for extreme teams
            league_confidence = 0.7 + 0.3 * (league_multiplier - 1.0) / 0.25  # Higher confidence for top leagues
            
            confidence = max(0.3, min(0.95, (sample_confidence + consistency_confidence + league_confidence) / 3.0))
            
            # Store enhanced team strengths with much wider, more meaningful ranges
            self.team_strengths[team] = {
                'attack': max(0.4, min(2.0, attack_strength)),
                'defense': max(0.4, min(2.0, defense_strength)), 
                'overall': max(0.5, min(1.8, overall_strength)),
                'form': max(0.5, min(1.8, form_strength)),
                'confidence': confidence,
                'league': team_league,
                'style': style,
                'home_strength': home_advantage,
                'momentum': max(-1.0, min(1.0, momentum)),
                
                # Additional detailed stats for better predictions
                'win_rate': win_rate,
                'draw_rate': draw_rate,
                'goals_per_game': avg_goals_for,
                'goals_conceded_per_game': avg_goals_against,
                'clean_sheet_rate': clean_sheet_rate,
                'home_goals_per_game': home_goals_per_game,
                'away_goals_per_game': away_goals_per_game,
                'home_win_rate': home_win_rate,
                'away_win_rate': away_win_rate,
                'quality_factor': max(0.6, min(1.6, quality_factor))
            }
        
        print(f"✅ Enhanced team strengths calculated for {len(self.team_strengths)} teams")
        
        # Print some examples for debugging
        sample_teams = list(self.team_strengths.keys())[:5]
        print("\n📊 Sample team strength analysis:")
        for team in sample_teams:
            stats = self.team_strengths[team]
            print(f"  {team}: Attack {stats['attack']:.2f} | Defense {stats['defense']:.2f} | "
                  f"Style: {stats['style']} | Win Rate: {stats['win_rate']:.1%} | "
                  f"Goals: {stats['goals_per_game']:.1f} | Confidence: {stats['confidence']:.1%}")
        
        # Store enhanced league averages
        for league in self.data['Div'].unique():
            league_data = self.data[self.data['Div'] == league]
            self.league_averages[league] = {
                'home_goals': league_data['FTHG'].mean(),
                'away_goals': league_data['FTAG'].mean(),
                'total_goals': (league_data['FTHG'] + league_data['FTAG']).mean(),
                'quality_factor': league_stats[league]['avg_total_goals'] / 2.7,  # Relative to average
                'competitiveness': 1.0 - abs(league_stats[league]['home_win_rate'] - 0.45)  # How even the league is
            }
        
        print(f"📊 Enhanced league baselines calculated for {len(self.league_averages)} leagues")

if __name__ == "__main__":
    # Test the predictor
    predictor = SoccerPredictor()
    predictor.load_data()
    predictor.train()
    
    # Test prediction
    result = predictor.predict_match("Premier League", "Liverpool", "Arsenal")
    print(f"\n🔮 Test Prediction: Liverpool vs Arsenal")
    print(f"Expected Score: {result['most_likely_home_score']}-{result['most_likely_away_score']}")
    print(f"Probabilities: Home {result['probabilities']['home_win']:.1%} | "
          f"Draw {result['probabilities']['draw']:.1%} | "
          f"Away {result['probabilities']['away_win']:.1%}") 