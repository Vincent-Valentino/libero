#!/usr/bin/env python3
"""
Simple script to start the ML service for testing
"""

import subprocess
import sys
import os

def main():
    print("🚀 Starting Football ML Prediction Service...")
    
    # Check if we're in the right directory
    if not os.path.exists('main.py'):
        print("❌ Error: main.py not found. Please run from the libero-ml directory.")
        sys.exit(1)
    
    # Check if requirements are installed
    try:
        import fastapi
        import uvicorn
        import pandas
        import sklearn
        import scipy
        print("✅ All dependencies found")
    except ImportError as e:
        print(f"❌ Missing dependency: {e}")
        print("Please run: pip install -r requirements.txt")
        sys.exit(1)
    
    # Start the service
    print("🔥 Starting FastAPI service on http://localhost:8001")
    print("📊 Model will train automatically on startup...")
    print("🔮 Ready for predictions!")
    print("=" * 50)
    
    try:
        subprocess.run([
            sys.executable, "-m", "uvicorn", 
            "main:app", 
            "--host", "0.0.0.0", 
            "--port", "8001", 
            "--reload"
        ], check=True)
    except KeyboardInterrupt:
        print("\n🛑 Service stopped by user")
    except subprocess.CalledProcessError as e:
        print(f"❌ Error starting service: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main() 