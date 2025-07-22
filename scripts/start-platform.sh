#!/bin/bash

# VideoTube Platform Startup Script
# This script starts the complete video processing platform

set -e

echo "🎬 Starting VideoTube Platform"
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

print_status "Checking existing containers..."
if docker-compose ps | grep -q "Up"; then
    print_warning "Some containers are already running. Stopping them first..."
    docker-compose down
fi

print_status "Building and starting all services..."
docker-compose up -d --build

print_status "Waiting for services to be healthy..."
sleep 30

# Check service health
print_status "Checking service health..."

# Check MongoDB
if docker-compose exec -T mongodb mongosh --eval "db.adminCommand('ismaster')" --quiet > /dev/null 2>&1; then
    print_success "MongoDB is healthy"
else
    print_warning "MongoDB might still be starting..."
fi

# Check Redis
if docker-compose exec -T redis redis-cli ping | grep -q "PONG"; then
    print_success "Redis is healthy"
else
    print_warning "Redis might still be starting..."
fi

# Check MinIO
if curl -s http://localhost:9000/minio/health/live > /dev/null; then
    print_success "MinIO is healthy"
else
    print_warning "MinIO might still be starting..."
fi

# Check Backend API
if curl -s http://localhost:8080/health > /dev/null; then
    print_success "Backend API is healthy"
else
    print_warning "Backend API might still be starting..."
fi

# Check Frontend
if curl -s http://localhost:3000 > /dev/null; then
    print_success "Frontend is healthy"
else
    print_warning "Frontend might still be starting..."
fi

echo ""
print_success "VideoTube Platform is starting up!"
echo ""
echo "📍 Access Points:"
echo "   🌐 Frontend:      http://localhost:3000"
echo "   🔧 Backend API:   http://localhost:8080"
echo "   📊 MinIO Console: http://localhost:9001 (minioadmin/minioadmin)"
echo ""
echo "🔍 Useful Commands:"
echo "   View logs:        docker-compose logs -f"
echo "   Stop platform:    docker-compose down"
echo "   Restart:          docker-compose restart"
echo ""
echo "🧪 Test Upload:"
echo "   Run test script:  ./scripts/test-upload.sh"
echo ""

# Show container status
print_status "Container Status:"
docker-compose ps

echo ""
print_success "Setup complete! Your video processing platform is ready to use."
echo ""
echo "💡 Tips:"
echo "   - It may take a few minutes for all services to be fully ready"
echo "   - Check logs if you encounter any issues: docker-compose logs [service-name]"
echo "   - The platform processes videos with FFmpeg in the background"
echo "" 