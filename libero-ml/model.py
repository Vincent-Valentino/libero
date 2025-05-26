"""
model.py: Comprehensive soccer match prediction using historical data from multiple leagues.

Dependencies:
- pandas 
- scikit-learn
- numpy
"""

import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, StandardScaler
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import accuracy_score, classification_report
import glob
import os
from datetime import datetime

class SoccerPredictor:
    def __init__(self):
        self.model = RandomForestClassifier(n_estimators=100, random_state=42)
        self.team_encoder = LabelEncoder()
        self.league_encoder = LabelEncoder()
        self.scaler = StandardScaler()
        
    def load_data(self, data_dir='data'):
        """Load and combine all historical match data from CSV files"""
        all_files = glob.glob(os.path.join(data_dir, '*.csv'))
        dataframes = []
        
        for file in all_files:
            df = pd.read_csv(file)
            # Convert Date to datetime
            df['Date'] = pd.to_datetime(df['Date'], format='%d/%m/%Y')
            dataframes.append(df)
            
        self.data = pd.concat(dataframes, ignore_index=True)
        # Sort by date to maintain temporal order
        self.data = self.data.sort_values('Date')
        return self.data
    
    def preprocess_data(self):
        """Process the raw data into features and target variables using all available metrics"""
        # Encode categorical variables
        self.data['HomeTeam_encoded'] = self.team_encoder.fit_transform(self.data['HomeTeam'])
        self.data['AwayTeam_encoded'] = self.team_encoder.transform(self.data['AwayTeam'])
        self.data['League_encoded'] = self.league_encoder.fit_transform(self.data['Div'])
        
        # Extract time features
        self.data['DayOfWeek'] = self.data['Date'].dt.dayofweek
        self.data['Month'] = self.data['Date'].dt.month
        
        # Create features including all available metrics
        features = [
            'League_encoded',
            'HomeTeam_encoded', 'AwayTeam_encoded',
            'DayOfWeek', 'Month',
            'HTHG', 'HTAG',  # Half-time goals
            'HS', 'AS',      # Shots
            'HST', 'AST',    # Shots on target
            'HF', 'AF',      # Fouls
            'HC', 'AC',      # Corners
            'HY', 'AY',      # Yellow cards
            'HR', 'AR'       # Red cards
        ]
        
        # Create target (1: Home win, 0: Draw, -1: Away win)
        self.data['Target'] = np.where(self.data['FTHG'] > self.data['FTAG'], 1,
                                     np.where(self.data['FTHG'] == self.data['FTAG'], 0, -1))
        
        # Prepare feature matrix
        X = self.data[features].copy()
        y = self.data['Target']
        
        # Scale numerical features
        numerical_features = [f for f in features if f not in 
                            ['League_encoded', 'HomeTeam_encoded', 'AwayTeam_encoded']]
        X[numerical_features] = self.scaler.fit_transform(X[numerical_features])
        
        return train_test_split(X, y, test_size=0.2, random_state=42)
    
    def train(self):
        """Train the model on the preprocessed data"""
        X_train, X_test, y_train, y_test = self.preprocess_data()
        self.model.fit(X_train, y_train)
        
        # Evaluate model
        train_pred = self.model.predict(X_train)
        test_pred = self.model.predict(X_test)
        
        print("Training Accuracy:", accuracy_score(y_train, train_pred))
        print("Test Accuracy:", accuracy_score(y_test, test_pred))
        print("\nClassification Report:")
        print(classification_report(y_test, test_pred))
        
        # Feature importance
        feature_names = [
            'League', 'Home Team', 'Away Team', 'Day of Week', 'Month',
            'Half-time Home Goals', 'Half-time Away Goals',
            'Home Shots', 'Away Shots',
            'Home Shots on Target', 'Away Shots on Target',
            'Home Fouls', 'Away Fouls',
            'Home Corners', 'Away Corners',
            'Home Yellow Cards', 'Away Yellow Cards',
            'Home Red Cards', 'Away Red Cards'
        ]
        
        feature_importance = pd.DataFrame({
            'feature': feature_names,
            'importance': self.model.feature_importances_
        })
        print("\nFeature Importance:")
        print(feature_importance.sort_values('importance', ascending=False))
        
    def predict_match(self, league, home_team, away_team, stats):
        """Predict the outcome of a match given all available statistics"""
        # Encode categorical inputs
        league_encoded = self.league_encoder.transform([league])[0]
        home_encoded = self.team_encoder.transform([home_team])[0]
        away_encoded = self.team_encoder.transform([away_team])[0]
        
        # Get current date features
        current_date = datetime.now()
        day_of_week = current_date.weekday()
        month = current_date.month
        
        # Prepare feature vector
        features = np.array([[
            league_encoded, home_encoded, away_encoded,
            day_of_week, month,
            stats['HTHG'], stats['HTAG'],
            stats['HS'], stats['AS'],
            stats['HST'], stats['AST'],
            stats['HF'], stats['AF'],
            stats['HC'], stats['AC'],
            stats['HY'], stats['AY'],
            stats['HR'], stats['AR']
        ]])
        
        # Scale numerical features
        features[:, 3:] = self.scaler.transform(features[:, 3:])
        
        # Make prediction
        prediction = self.model.predict(features)[0]
        probabilities = self.model.predict_proba(features)[0]
        
        return {
            'prediction': prediction,  # 1: Home win, 0: Draw, -1: Away win
            'probabilities': {
                'home_win': probabilities[2],
                'draw': probabilities[1],
                'away_win': probabilities[0]
            }
        }

if __name__ == "__main__":
    predictor = SoccerPredictor()
    data = predictor.load_data()
    predictor.train()
