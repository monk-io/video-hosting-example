#!/bin/bash

# Test script for video upload and processing

API_URL="http://localhost:8080"

echo "🎬 Testing Video Processing Platform"
echo "====================================="

# Check if API is running
echo "🔍 Checking API health..."
if curl -s "$API_URL/health" > /dev/null; then
    echo "✅ API is running"
else
    echo "❌ API is not responding. Please start the services with:"
    echo "   docker-compose up -d"
    exit 1
fi

# Create a test video file (if it doesn't exist)
TEST_VIDEO="test-video.mp4"
if [ ! -f "$TEST_VIDEO" ]; then
    echo "📹 Creating test video file..."
    # Create a 10-second test video using FFmpeg
    if command -v ffmpeg &> /dev/null; then
        ffmpeg -f lavfi -i testsrc=duration=10:size=320x240:rate=30 -f lavfi -i sine=frequency=1000:duration=10 -c:v libx264 -c:a aac -shortest "$TEST_VIDEO" -y
        echo "✅ Test video created: $TEST_VIDEO"
    else
        echo "❌ FFmpeg not found. Please create a test video file named '$TEST_VIDEO'"
        echo "   Or install FFmpeg to generate one automatically"
        exit 1
    fi
fi

# Upload the test video
echo ""
echo "⬆️  Uploading test video..."
UPLOAD_RESPONSE=$(curl -s -X POST \
    -F "video=@$TEST_VIDEO" \
    -F "title=Test Video - $(date)" \
    -F "description=Automated test upload" \
    -F "uploaded_by=test-user" \
    "$API_URL/api/v1/videos/upload")

echo "Upload response: $UPLOAD_RESPONSE"

# Extract video ID from response
VIDEO_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"video_id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$VIDEO_ID" ]; then
    echo "❌ Failed to get video ID from upload response"
    exit 1
fi

echo "✅ Video uploaded successfully!"
echo "📝 Video ID: $VIDEO_ID"

# Check video details
echo ""
echo "📋 Checking video details..."
curl -s "$API_URL/api/v1/videos/$VIDEO_ID" | jq '.' || echo "Video details retrieved"

# Monitor processing jobs
echo ""
echo "🔄 Monitoring processing jobs..."
for i in {1..30}; do
    echo "Check $i/30..."
    
    # Get jobs for this video
    JOBS_RESPONSE=$(curl -s "$API_URL/api/v1/jobs/video/$VIDEO_ID")
    echo "Jobs status: $JOBS_RESPONSE" | jq '.' || echo "$JOBS_RESPONSE"
    
    # Check if all jobs are completed
    COMPLETED_COUNT=$(echo "$JOBS_RESPONSE" | grep -o '"status":"completed"' | wc -l)
    FAILED_COUNT=$(echo "$JOBS_RESPONSE" | grep -o '"status":"failed"' | wc -l)
    TOTAL_JOBS=$(echo "$JOBS_RESPONSE" | grep -o '"status":' | wc -l)
    
    if [ "$COMPLETED_COUNT" -eq 4 ]; then
        echo "✅ All jobs completed successfully!"
        break
    elif [ "$FAILED_COUNT" -gt 0 ]; then
        echo "❌ Some jobs failed"
        break
    fi
    
    sleep 5
done

# Final video status
echo ""
echo "📊 Final video status:"
curl -s "$API_URL/api/v1/videos/$VIDEO_ID" | jq '.' || curl -s "$API_URL/api/v1/videos/$VIDEO_ID"

echo ""
echo "🎉 Test completed!"
echo ""
echo "📍 Useful commands:"
echo "   View all videos: curl $API_URL/api/v1/videos"
echo "   Stream video: curl $API_URL/api/v1/videos/$VIDEO_ID/stream"
echo "   View jobs: curl $API_URL/api/v1/jobs/video/$VIDEO_ID"
echo ""
echo "🌐 Access MinIO console: http://localhost:9001 (minioadmin/minioadmin)" 