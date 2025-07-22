#!/bin/bash

echo "🎬 Setting up Video Processing Platform..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. You may want to review and modify it."
fi

# Build and start services
echo "🚀 Building and starting services..."
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 30

# Check service health
echo "🔍 Checking service health..."

# Check MongoDB
echo "Checking MongoDB..."
docker exec youtube_mongodb mongosh --quiet --eval "db.adminCommand('ismaster')" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ MongoDB is ready"
else
    echo "❌ MongoDB is not ready"
fi

# Check Redis
echo "Checking Redis..."
docker exec youtube_redis redis-cli ping > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Redis is ready"
else
    echo "❌ Redis is not ready"
fi

# Check MinIO
echo "Checking MinIO..."
curl -s http://localhost:9000/minio/health/live > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ MinIO is ready"
else
    echo "❌ MinIO is not ready"
fi

# Check Backend API
echo "Checking Backend API..."
sleep 5
curl -s http://localhost:8080/health > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Backend API is ready"
else
    echo "❌ Backend API is not ready"
fi

echo ""
echo "🎉 Setup completed!"
echo ""
echo "📍 Service URLs:"
echo "   - Frontend: http://localhost:3000"
echo "   - Backend API: http://localhost:8080"
echo "   - MinIO Console: http://localhost:9001 (minioadmin/minioadmin)"
echo ""
echo "🧪 Test the API:"
echo "   curl http://localhost:8080/health"
echo ""
echo "📋 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down" 