#!/bin/bash

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
sleep 10

# MinIO client configuration
mc alias set myminio http://minio:9000 minioadmin minioadmin

# Create buckets
echo "Creating buckets..."
mc mb myminio/videos --ignore-existing
mc mb myminio/thumbnails --ignore-existing

# Set bucket policies for public read access to processed videos and thumbnails
echo "Setting bucket policies..."
mc anonymous set public myminio/videos/processed/
mc anonymous set public myminio/videos/thumbnails/
mc anonymous set public myminio/thumbnails/

echo "MinIO initialization completed!" 