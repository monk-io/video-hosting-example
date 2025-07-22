# 🎬 Video Processing & Sharing Platform

A complete YouTube-like video processing and sharing platform built with modern technologies, featuring **parallel video processing workers**, **FFmpeg transcoding**, **thumbnail generation**, and containerized deployment.

## 🏗️ Architecture

This platform follows a microservices architecture with **distributed video processing**:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React Frontend │    │   Go Backend API │    │  Video Workers  │
│   (Port: 3000)  │◄──►│   (Port: 8080)   │◄──►│  (Parallel)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         │                        ▼                        │
         │              ┌─────────────────┐                │
         │              │     MongoDB     │                │
         │              │  (Port: 27017)  │                │
         │              └─────────────────┘                │
         │                        │                        │
         │                        ▼                        │
         │              ┌─────────────────┐                │
         └─────────────►│   Redis Queue   │◄───────────────┘
                        │   (Port: 6379)  │
                        └─────────────────┘
                                 │
                                 ▼
                       ┌─────────────────┐
                       │      MinIO      │
                       │   (Port: 9000)  │
                       └─────────────────┘
```

### Technology Stack
- **Frontend**: React 18 + TypeScript + Tailwind CSS
- **Backend**: Go 1.21 + Gin + Domain-Driven Design (DDD)
- **Workers**: Go + FFmpeg for video processing
- **Database**: MongoDB 7.0 (metadata, user data)
- **Queue**: Redis 7.0 (distributed job processing)
- **Storage**: MinIO (S3-compatible, video files)
- **Processing**: FFmpeg (video transcoding & thumbnails)
- **Deployment**: Docker Compose

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd youtube-example
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings if needed
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Wait for initialization (about 30-60 seconds)**
   ```bash
   # Check service status
   docker-compose ps
   
   # View logs
   docker-compose logs -f
   ```

5. **Access the application**
   - **Backend API**: http://localhost:8080
   - **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
   - **Frontend**: http://localhost:3000 (when implemented)

## 🎯 Features

### ✅ **Fully Implemented**

**🎥 Video Processing Pipeline**
- ✅ **Video Upload**: Multipart upload with validation
- ✅ **Distributed Processing**: Redis job queue with parallel workers
- ✅ **FFmpeg Transcoding**: 480p, 720p, 1080p quality options
- ✅ **Thumbnail Generation**: Automatic thumbnail creation
- ✅ **Progress Tracking**: Real-time job status monitoring
- ✅ **File Storage**: Organized MinIO structure (original/processed/thumbnails)

**🔧 Backend Infrastructure**
- ✅ **RESTful API**: Complete CRUD operations
- ✅ **Domain-Driven Design**: Clean architecture
- ✅ **Job Queue**: Redis-based distributed processing
- ✅ **Database Operations**: MongoDB with proper indexing
- ✅ **File Streaming**: Direct video streaming with quality selection

**🏗️ DevOps & Deployment**
- ✅ **Docker Compose**: Complete containerized stack
- ✅ **Auto-initialization**: MinIO buckets and permissions
- ✅ **Health Checks**: Service monitoring endpoints
- ✅ **Structured Logging**: Comprehensive error tracking

### 🚧 **Next Phase** (Ready for Implementation)
- [ ] **React Frontend**: Upload interface and video player
- [ ] **User Authentication**: JWT-based auth system
- [ ] **Real-time Updates**: WebSocket progress notifications
- [ ] **Advanced Features**: Video search, playlists, comments

## 🛠️ API Endpoints

### Videos
```bash
# Upload video
POST   /api/v1/videos/upload
curl -X POST -F "video=@video.mp4" -F "title=My Video" http://localhost:8080/api/v1/videos/upload

# List videos (paginated)
GET    /api/v1/videos?page=1&limit=20
curl http://localhost:8080/api/v1/videos

# Get video details
GET    /api/v1/videos/:id
curl http://localhost:8080/api/v1/videos/64a7b8c9d1e2f3a4b5c6d7e8

# Stream video (original or processed)
GET    /api/v1/videos/:id/stream?quality=720p
curl http://localhost:8080/api/v1/videos/64a7b8c9d1e2f3a4b5c6d7e8/stream

# Trigger manual processing
POST   /api/v1/videos/:id/process
curl -X POST http://localhost:8080/api/v1/videos/64a7b8c9d1e2f3a4b5c6d7e8/process
```

### Jobs
```bash
# Get job status
GET    /api/v1/jobs/:id
curl http://localhost:8080/api/v1/jobs/64a7b8c9d1e2f3a4b5c6d7e9

# Get all jobs for a video
GET    /api/v1/jobs/video/:videoId
curl http://localhost:8080/api/v1/jobs/video/64a7b8c9d1e2f3a4b5c6d7e8

# Get active processing jobs
GET    /api/v1/jobs/active
curl http://localhost:8080/api/v1/jobs/active
```

### Health
```bash
# Service health check
GET    /health
curl http://localhost:8080/health
```

## 🧪 Testing

### Automated Test
Run the complete test pipeline:
```bash
./scripts/test-upload.sh
```

This script will:
1. Create a test video (if FFmpeg is available)
2. Upload the video via API
3. Monitor processing jobs
4. Show final results

### Manual Testing
```bash
# 1. Check API health
curl http://localhost:8080/health

# 2. Upload a video
curl -X POST \
  -F "video=@your-video.mp4" \
  -F "title=Test Video" \
  -F "description=Test upload" \
  http://localhost:8080/api/v1/videos/upload

# 3. Monitor processing
curl http://localhost:8080/api/v1/jobs/active

# 4. Check final video status
curl http://localhost:8080/api/v1/videos/{VIDEO_ID}
```

## 📊 Video Processing Flow

1. **Upload** → Video uploaded to MinIO (`videos/original/`)
2. **Job Creation** → 4 jobs created (1 thumbnail + 3 transcode jobs)
3. **Queue Distribution** → Jobs sent to Redis queue
4. **Worker Processing** → 2 parallel workers process jobs
5. **FFmpeg Processing** → Videos transcoded to multiple formats
6. **Storage** → Processed files saved (`videos/processed/`, `thumbnails/`)
7. **Database Update** → Video metadata updated with new formats
8. **Completion** → Video status changed to "ready"

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Backend server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://admin:password@mongodb:27017/youtube?authSource=admin` |
| `REDIS_URI` | Redis connection string | `redis://redis:6379` |
| `MINIO_ENDPOINT` | MinIO server endpoint | `minio:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` |
| `FRONTEND_URL` | Frontend URL for CORS | `http://localhost:3000` |
| `WORKER_ID` | Unique worker identifier | Auto-generated |

### Video Processing Settings

- **Supported Input Formats**: MP4, AVI, MOV, WMV, FLV, WebM, MKV
- **Output Formats**: 
  - 480p: H.264, 1Mbps max bitrate
  - 720p: H.264, 2.5Mbps max bitrate  
  - 1080p: H.264, 4.5Mbps max bitrate
- **Audio**: AAC, 128kbps
- **Max File Size**: 1GB
- **Parallel Workers**: 2 (configurable)

## 📁 Project Structure

```
youtube-example/
├── backend/                 # Go backend API (22 files)
│   ├── internal/
│   │   ├── domain/         # Business logic (DDD)
│   │   │   ├── entities/   # Video, Job, User entities
│   │   │   ├── repositories/ # Data access interfaces
│   │   │   └── services/   # Business services
│   │   ├── infrastructure/ # External dependencies
│   │   │   ├── database/   # MongoDB client
│   │   │   ├── storage/    # MinIO client
│   │   │   ├── queue/      # Redis client & job publisher
│   │   │   └── repositories/ # Repository implementations
│   │   ├── application/    # Application services
│   │   │   └── handlers/   # HTTP handlers
│   │   └── interfaces/     # API layer
│   │       ├── http/       # Routes and middleware
│   │       └── middleware/ # CORS, logging
│   └── pkg/               # Shared packages
├── worker/                 # Video processing workers (10 files)
│   ├── internal/
│   │   ├── processor/     # FFmpeg transcoding logic
│   │   ├── storage/       # MinIO client
│   │   └── queue/         # Redis & MongoDB clients
│   └── pkg/              # Worker configuration
├── frontend/              # React frontend (future)
├── scripts/               # Setup and testing scripts
│   ├── setup.sh          # Complete platform setup
│   ├── test-upload.sh    # Video upload testing
│   └── init-minio.sh     # MinIO initialization
├── docker-compose.yml    # Container orchestration
└── .env.example         # Environment configuration
```

## 🏃‍♂️ Development

### Backend Development
```bash
cd backend
go mod tidy
go run main.go
```

### Worker Development
```bash
cd worker
go mod tidy
go run main.go
```

### Scaling Workers
```bash
# Scale to 4 workers
docker-compose up --scale worker1=2 --scale worker2=2
```

## 📊 Monitoring

### Container Status
```bash
docker-compose ps
```

### Logs
```bash
# All services
docker-compose logs -f

# Specific services
docker-compose logs -f backend worker1 worker2
```

### MinIO Console
- URL: http://localhost:9001
- Credentials: minioadmin/minioadmin
- View uploaded videos and processed files

### Database Access
```bash
# MongoDB shell
docker exec -it youtube_mongodb mongosh --username admin --password password --authenticationDatabase admin

# Redis CLI
docker exec -it youtube_redis redis-cli
```

## 🔒 Security Considerations

- ✅ File validation and size limits
- ✅ CORS configuration
- ✅ Error handling and logging
- 🚧 Authentication system (planned)
- 🚧 Rate limiting (planned)
- 🚧 File encryption (planned)

## 🚀 Production Deployment

### Recommended Setup
1. **Container Orchestration**: Kubernetes or Docker Swarm
2. **Load Balancer**: Nginx or HAProxy
3. **Database**: MongoDB Atlas or self-hosted cluster
4. **Storage**: AWS S3 or distributed MinIO
5. **Monitoring**: Prometheus + Grafana
6. **Logging**: ELK Stack

### Performance Optimizations
- Horizontal scaling of workers
- CDN for video delivery
- Database indexing
- Connection pooling
- Caching strategies

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

---

## 🎉 Current Status: **Production-Ready Backend & Workers**

✅ **Complete video processing pipeline with distributed workers**  
✅ **FFmpeg transcoding and thumbnail generation**  
✅ **RESTful API with job monitoring**  
✅ **Containerized deployment with Docker Compose**  
🚧 **React frontend ready for implementation**  

**Built with ❤️ using Go, FFmpeg, MongoDB, Redis, and MinIO.** 