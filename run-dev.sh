#!/bin/bash

echo "Starting backend..."
# Run backend in the background, redirecting output might be needed if noisy
(cd libero-backend && go run main.go app.go) &
BACKEND_PID=$!
echo "Backend started with PID $BACKEND_PID"

echo "Starting frontend..."
# Run frontend in the background
(cd libero-frontend && pnpm run dev) &
FRONTEND_PID=$!
echo "Frontend started with PID $FRONTEND_PID"

echo "Both processes are running in the background."
echo "You might need to manually stop them (e.g., using 'kill $BACKEND_PID $FRONTEND_PID' or Ctrl+C multiple times)."

# Keep the script running so the terminal doesn't close immediately (optional)
# wait