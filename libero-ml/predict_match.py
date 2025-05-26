from model import SoccerPredictor

def predict_upcoming_match():
    # Initialize and train the model
    predictor = SoccerPredictor()
    print("Loading and training model...")
    predictor.load_data()
    predictor.train()
    
    # Get match details from user
    print("\nEnter match details:")
    league = input("League (e.g., Premier League, La Liga, etc.): ")
    home_team = input("Home Team: ")
    away_team = input("Away Team: ")
    
    # Get match statistics (using average values if not available)
    print("\nEnter match statistics (press Enter to use average values):")
    stats = {
        'HTHG': input("Half-time Home Goals (default=0): ") or "0",
        'HTAG': input("Half-time Away Goals (default=0): ") or "0",
        'HS': input("Home Shots (default=12): ") or "12",
        'AS': input("Away Shots (default=10): ") or "10",
        'HST': input("Home Shots on Target (default=5): ") or "5",
        'AST': input("Away Shots on Target (default=4): ") or "4",
        'HF': input("Home Fouls (default=10): ") or "10",
        'AF': input("Away Fouls (default=10): ") or "10",
        'HC': input("Home Corners (default=5): ") or "5",
        'AC': input("Away Corners (default=4): ") or "4",
        'HY': input("Home Yellow Cards (default=2): ") or "2",
        'AY': input("Away Yellow Cards (default=2): ") or "2",
        'HR': input("Home Red Cards (default=0): ") or "0",
        'AR': input("Away Red Cards (default=0): ") or "0"
    }
    
    # Convert stats to integers
    stats = {k: int(v) for k, v in stats.items()}
    
    # Get prediction
    result = predictor.predict_match(league, home_team, away_team, stats)
    
    # Display results
    print("\nPrediction Results:")
    print("-" * 50)
    print(f"Match: {home_team} vs {away_team}")
    print(f"League: {league}")
    
    # Interpret prediction
    prediction = result['prediction']
    if prediction == 1:
        winner = home_team
    elif prediction == -1:
        winner = away_team
    else:
        winner = "Draw predicted"
    
    print(f"\nPredicted Winner: {winner}")
    print("\nProbabilities:")
    print(f"Home Win ({home_team}): {result['probabilities']['home_win']:.2%}")
    print(f"Draw: {result['probabilities']['draw']:.2%}")
    print(f"Away Win ({away_team}): {result['probabilities']['away_win']:.2%}")

if __name__ == "__main__":
    predict_upcoming_match() 