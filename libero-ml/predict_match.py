"""
Match prediction interface for the Poisson-based football score predictor.
Handles the conversion between the ML model output and API response format.
"""

from model import SoccerPredictor

def predict_match_result(predictor, league, home_team, away_team, stats=None):
    """
    Predict the outcome of a football match using the trained Poisson models
    
    Args:
        predictor: Trained SoccerPredictor instance with Poisson models
        league: League name (e.g., "Premier League", "La Liga")
        home_team: Home team name
        away_team: Away team name
        stats: Optional match statistics dict (team form, etc.)
        
    Returns:
        dict: Prediction results with exact scores, expected goals, and probabilities
    """
    
    # Get prediction from the Poisson models
    result = predictor.predict_match(league, home_team, away_team, stats)
    
    # Extract the core prediction data
    expected_home_goals = result['expected_home_goals']
    expected_away_goals = result['expected_away_goals']
    most_likely_home_score = result['most_likely_home_score']
    most_likely_away_score = result['most_likely_away_score']
    
    # Get match outcome probabilities
    probabilities = result['probabilities']
    home_win_prob = probabilities['home_win']
    draw_prob = probabilities['draw']
    away_win_prob = probabilities['away_win']
    
    # Determine the prediction result (1 = home win, 0 = draw, -1 = away win)
    if home_win_prob > max(draw_prob, away_win_prob):
        prediction = 1  # Home win
    elif away_win_prob > draw_prob:
        prediction = -1  # Away win
    else:
        prediction = 0  # Draw
    
    # Return structured result for API (matching FastAPI Pydantic model)
    return {
        'prediction': prediction,
        'probabilities': probabilities,
        'expected_home_goals': expected_home_goals,
        'expected_away_goals': expected_away_goals,
        'most_likely_home_score': most_likely_home_score,
        'most_likely_away_score': most_likely_away_score
    }


def predict_upcoming_match():
    """
    Interactive function for manual match prediction using Poisson models
    """
    # Initialize and train the Poisson models
    predictor = SoccerPredictor()
    print("🚀 Loading and training Poisson models...")
    predictor.load_data()
    predictor.train()
    
    print("\n" + "="*60)
    print("🔮 FOOTBALL SCORE PREDICTOR")
    print("="*60)
    
    # Get match details from user
    print("\nEnter match details:")
    league = input("League (e.g., Premier League, La Liga, etc.): ").strip()
    home_team = input("Home Team: ").strip()
    away_team = input("Away Team: ").strip()
    
    # Optional: Get team form data
    print(f"\nOptional: Enter recent form data (press Enter to use defaults)")
    print("Recent form helps improve prediction accuracy...")
    
    advanced_input = input("Would you like to enter team form data? (y/n): ").strip().lower()
    
    stats = None
    if advanced_input == 'y':
        try:
            print(f"\nRecent form for {home_team} (at home):")
            home_goals_avg = input("Average goals scored in last 5 home games (default=1.5): ").strip()
            home_conceded_avg = input("Average goals conceded in last 5 home games (default=1.2): ").strip()
            home_form_points = input("Points per game in last 5 home games (default=1.5): ").strip()
            
            print(f"\nRecent form for {away_team} (away from home):")
            away_goals_avg = input("Average goals scored in last 5 away games (default=1.3): ").strip()
            away_conceded_avg = input("Average goals conceded in last 5 away games (default=1.4): ").strip()
            away_form_points = input("Points per game in last 5 away games (default=1.3): ").strip()
            
            # Create stats dictionary
            stats = {}
            if home_goals_avg:
                stats['home_goals_avg'] = float(home_goals_avg)
            if home_conceded_avg:
                stats['home_conceded_avg'] = float(home_conceded_avg)
            if home_form_points:
                stats['home_form_points'] = float(home_form_points)
            if away_goals_avg:
                stats['away_goals_avg'] = float(away_goals_avg)
            if away_conceded_avg:
                stats['away_conceded_avg'] = float(away_conceded_avg)
            if away_form_points:
                stats['away_form_points'] = float(away_form_points)
                
        except ValueError:
            print("⚠️  Invalid input, using default values...")
            stats = None
    
    # Get prediction using the Poisson models
    print(f"\n🔍 Analyzing match: {home_team} vs {away_team}")
    print("⚡ Running Poisson regression models...")
    
    result = predict_match_result(predictor, league, home_team, away_team, stats)
    
    # Display results in a nice format
    print("\n" + "="*60)
    print("🎯 PREDICTION RESULTS")
    print("="*60)
    
    print(f"🏟️  Match: {home_team} vs {away_team}")
    print(f"🏆 League: {league}")
    print(f"📅 Prediction Date: {predictor.data['Date'].max().strftime('%Y-%m-%d') if predictor.data is not None else 'N/A'}")
    
    print(f"\n📊 EXACT SCORE PREDICTION:")
    print(f"🎯 Most Likely Score: {result['most_likely_home_score']} - {result['most_likely_away_score']}")
    print(f"⚽ Expected Goals: {result['expected_home_goals']:.2f} - {result['expected_away_goals']:.2f}")
    print(f"🏅 Result: {result['prediction']}")
    
    print(f"\n📈 MATCH OUTCOME PROBABILITIES:")
    print(f"🏠 {home_team} Win: {result['probabilities']['home_win']:.1%}")
    print(f"🤝 Draw: {result['probabilities']['draw']:.1%}")
    print(f"✈️  {away_team} Win: {result['probabilities']['away_win']:.1%}")
    
    # Additional insights
    print(f"\n💡 INSIGHTS:")
    total_expected_goals = result['expected_home_goals'] + result['expected_away_goals']
    if total_expected_goals > 2.5:
        print(f"📊 High-scoring match expected ({total_expected_goals:.1f} total goals)")
        print(f"⭐ Recommendation: Over 2.5 goals")
    else:
        print(f"📊 Low-scoring match expected ({total_expected_goals:.1f} total goals)")
        print(f"⭐ Recommendation: Under 2.5 goals")
    
    if abs(result['expected_home_goals'] - result['expected_away_goals']) < 0.3:
        print(f"⚖️  Very close match - small goal difference expected")
    elif result['expected_home_goals'] > result['expected_away_goals'] + 0.5:
        print(f"🏠 Home team has clear advantage")
    elif result['expected_away_goals'] > result['expected_home_goals'] + 0.5:
        print(f"✈️  Away team has clear advantage")
    
    print("\n" + "="*60)
    print("🔬 Powered by Poisson Regression Models")
    print("📊 Based on team form, betting odds, and historical data")
    print("="*60)


if __name__ == "__main__":
    predict_upcoming_match() 