"""
FastAPI microservice for football match prediction using Poisson regression models.
Provides exact score predictions based on team form and historical data.
"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional, Dict, Any
import uvicorn
import os
import logging
from datetime import datetime

# Import our Poisson-based prediction modules
from model import SoccerPredictor
from predict_match import predict_match_result

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global variable to store the trained predictor
predictor = None

# FastAPI app initialization
app = FastAPI(
    title="Football Score Predictor API",
    description="Poisson-based football score prediction with exact scores and probabilities",
    version="1.0.0"
)

# Add CORS middleware for cross-origin requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify your frontend domain
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Request/Response models
class PredictionRequest(BaseModel):
    league: str
    home_team: str
    away_team: str
    stats: Optional[Dict[str, Any]] = None

class PredictionResponse(BaseModel):
    prediction: int
    probabilities: Dict[str, float]
    expected_home_goals: float
    expected_away_goals: float
    most_likely_home_score: int
    most_likely_away_score: int

# Startup event to load and train model
@app.on_event("startup")
async def startup_event():
    global predictor
    try:
        logger.info("🚀 Loading and training Poisson-based football predictor...")
        
        # Initialize and train the predictor
        predictor = SoccerPredictor()
        predictor.load_data()
        predictor.train()
        
        logger.info("✅ Model training completed successfully!")
        
        # Run debug and test functions
        debug_model_features()
        test_predictions()
        
    except Exception as e:
        logger.error(f"❌ Failed to initialize predictor: {e}")
        import traceback
        traceback.print_exc()
        raise e

@app.get("/")
async def root():
    """Health check endpoint"""
    return {
        "message": "⚽ Football Score Predictor API",
        "status": "running",
        "model_ready": predictor is not None
    }

@app.post("/predict", response_model=PredictionResponse)
async def predict_match(request: PredictionRequest):
    """
    Predict football match outcome using Poisson regression models
    """
    global predictor
    
    if predictor is None:
        raise HTTPException(
            status_code=503, 
            detail="Model not ready. Please wait for training to complete."
        )
    
    try:
        # Make prediction using the trained model
        result = predict_match_result(
            predictor=predictor,
            league=request.league,
            home_team=request.home_team,
            away_team=request.away_team,
            stats=request.stats
        )
        
        logger.info(f"🔮 Prediction made: {request.home_team} vs {request.away_team}")
        return PredictionResponse(**result)
        
    except Exception as e:
        logger.error(f"❌ Prediction error: {e}")
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")

@app.get("/teams")
async def get_available_teams():
    """
    Get list of all available teams in the dataset
    """
    global predictor
    
    if predictor is None:
        raise HTTPException(
            status_code=503, 
            detail="Model not ready. Please wait for training to complete."
        )
    
    try:
        teams = predictor.get_available_teams()
        return {"teams": teams}
    except Exception as e:
        logger.error(f"❌ Error getting teams: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to get teams: {str(e)}")

@app.get("/leagues")
async def get_available_leagues():
    """
    Get list of all available leagues in the dataset
    """
    global predictor
    
    if predictor is None:
        raise HTTPException(
            status_code=503, 
            detail="Model not ready. Please wait for training to complete."
        )
    
    try:
        leagues = predictor.get_available_leagues()
        return {"leagues": leagues}
    except Exception as e:
        logger.error(f"❌ Error getting leagues: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to get leagues: {str(e)}")

@app.get("/health")
async def health_check():
    """
    Detailed health check with model status
    """
    return {
        "status": "healthy",
        "model_loaded": predictor is not None,
        "timestamp": datetime.now().isoformat(),
        "service": "football-predictor-ml"
    }

def test_predictions():
    """Test function to debug prediction results"""
    print("\n" + "="*60)
    print("🧪 TESTING PREDICTION SYSTEM")
    print("="*60)
    
    # Test matchups
    test_cases = [
        ("D1", "Bayern Munich", "Borussia Dortmund"),
        ("SP1", "Barcelona", "Alaves"), 
        ("SP1", "Alaves", "Barcelona"),
        ("E0", "Arsenal", "Angers"),
        ("E0", "Bayern Munich", "Arsenal")  # Cross-league
    ]
    
    for league, home_team, away_team in test_cases:
        print(f"\n🔮 Testing: {home_team} vs {away_team} ({league})")
        print("-" * 50)
        
        try:
            # Make prediction
            result = predictor.predict_match(league, home_team, away_team)
            
            # Print detailed results
            print(f"📊 Expected Goals: {result['expected_home_goals']:.2f} - {result['expected_away_goals']:.2f}")
            print(f"🎯 Most Likely Score: {result['most_likely_home_score']} - {result['most_likely_away_score']}")
            print(f"📈 Outcome Probabilities:")
            print(f"   🏠 {home_team} Win: {result['probabilities']['home_win']:.1%}")
            print(f"   🤝 Draw: {result['probabilities']['draw']:.1%}")
            print(f"   ✈️  {away_team} Win: {result['probabilities']['away_win']:.1%}")
            
            # Determine predicted outcome
            if result['prediction'] == 1:
                outcome = f"🏠 {home_team} Win"
            elif result['prediction'] == -1:
                outcome = f"✈️ {away_team} Win"
            else:
                outcome = "🤝 Draw"
            print(f"🏆 Predicted Outcome: {outcome}")
            
        except Exception as e:
            print(f"❌ Error predicting {home_team} vs {away_team}: {e}")
            import traceback
            traceback.print_exc()
    
    print("\n" + "="*60)
    print("🧪 TESTING COMPLETED")
    print("="*60)

def debug_model_features():
    """Debug function to show model features and data"""
    print("\n" + "="*60)
    print("🔍 MODEL DEBUG INFORMATION")
    print("="*60)
    
    print(f"\n📊 Dataset Info:")
    print(f"   Total matches: {len(predictor.data)}")
    print(f"   Date range: {predictor.data['Date'].min()} to {predictor.data['Date'].max()}")
    print(f"   Leagues: {', '.join(predictor.data['Div'].unique())}")
    
    print(f"\n🏠 Home Features ({len(predictor.home_features)}):")
    for i, feature in enumerate(predictor.home_features):
        print(f"   {i+1:2d}. {feature}")
    
    print(f"\n✈️  Away Features ({len(predictor.away_features)}):")
    for i, feature in enumerate(predictor.away_features):
        print(f"   {i+1:2d}. {feature}")
    
    # Check for enhanced features
    enhanced_features = ['market_home_prob', 'shot_ratio', 'corner_ratio', 'foul_ratio']
    print(f"\n🚀 Enhanced Features Status:")
    for feature in enhanced_features:
        status = "✅" if feature in predictor.data.columns else "❌"
        print(f"   {status} {feature}")
    
    # Sample some recent data
    print(f"\n📈 Sample Recent Data (last 5 matches):")
    recent_data = predictor.data.tail(5)[['Date', 'HomeTeam', 'AwayTeam', 'FTHG', 'FTAG']]
    for _, match in recent_data.iterrows():
        print(f"   {match['Date'].strftime('%Y-%m-%d')}: {match['HomeTeam']} {match['FTHG']}-{match['FTAG']} {match['AwayTeam']}")
    
    print("\n" + "="*60)

if __name__ == "__main__":
    # Start the server - model initialization happens in startup_event
    logger.info("🚀 Starting Football Prediction API...")
    logger.info("📈 Model will be trained automatically on startup...")
    
    # Start the server
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8001,
        log_level="info"
    )