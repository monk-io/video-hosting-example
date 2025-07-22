import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { ArrowLeftIcon, ShareIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { useVideo, useJobs } from '../hooks/useVideo';
import VideoPlayer from '../components/VideoPlayer';
import JobStatus from '../components/JobStatus';

const VideoDetailPage: React.FC = () => {
  const { videoId } = useParams<{ videoId: string }>();
  const { video, loading: videoLoading, error: videoError } = useVideo(videoId);
  const { jobs, loading: jobsLoading } = useJobs(videoId);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  };

  const formatDuration = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  };

  const handleShare = async () => {
    if (navigator.share && video) {
      try {
        await navigator.share({
          title: video.title,
          text: video.description,
          url: window.location.href,
        });
      } catch (err) {
        // Fallback to copying URL
        navigator.clipboard.writeText(window.location.href);
        alert('Link copied to clipboard!');
      }
    } else {
      // Fallback to copying URL
      navigator.clipboard.writeText(window.location.href);
      alert('Link copied to clipboard!');
    }
  };

  if (videoLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="animate-pulse">
            <div className="aspect-w-16 aspect-h-9 bg-gray-300 rounded-lg mb-6"></div>
            <div className="h-8 bg-gray-300 rounded w-3/4 mb-4"></div>
            <div className="h-4 bg-gray-300 rounded w-1/4 mb-6"></div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              <div className="lg:col-span-2">
                <div className="h-32 bg-gray-300 rounded"></div>
              </div>
              <div className="h-64 bg-gray-300 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (videoError || !video) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Video Not Found</h2>
          <p className="text-gray-600 mb-6">
            {videoError || 'The video you\'re looking for doesn\'t exist.'}
          </p>
          <Link to="/" className="btn-primary">
            <ArrowLeftIcon className="h-5 w-5 mr-2" />
            Back to Home
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="flex items-center text-gray-600 hover:text-gray-900">
              <ArrowLeftIcon className="h-5 w-5 mr-2" />
              Back to Videos
            </Link>
            <h1 className="text-xl font-bold text-youtube-red">VideoTube</h1>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Video Player */}
        <div className="mb-8">
          <div className="video-container bg-black rounded-lg overflow-hidden">
            {video.status === 'ready' ? (
              <VideoPlayer video={video} autoPlay={false} />
            ) : (
              <div className="flex items-center justify-center h-full">
                <div className="text-center text-white">
                  {video.status === 'processing' && (
                    <>
                      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
                      <h3 className="text-lg font-medium mb-2">Processing Video...</h3>
                      <p className="text-gray-300">Your video is being transcoded and will be ready shortly.</p>
                    </>
                  )}
                  {video.status === 'uploaded' && (
                    <>
                      <div className="h-12 w-12 bg-gray-600 rounded-full mx-auto mb-4 flex items-center justify-center">
                        <span className="text-lg">⏳</span>
                      </div>
                      <h3 className="text-lg font-medium mb-2">Processing Queued</h3>
                      <p className="text-gray-300">Your video is in the processing queue.</p>
                    </>
                  )}
                  {video.status === 'failed' && (
                    <>
                      <div className="h-12 w-12 bg-red-600 rounded-full mx-auto mb-4 flex items-center justify-center">
                        <span className="text-lg">❌</span>
                      </div>
                      <h3 className="text-lg font-medium mb-2">Processing Failed</h3>
                      <p className="text-gray-300">There was an error processing your video.</p>
                    </>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Video Info */}
          <div className="lg:col-span-2">
            <div className="card p-6">
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <h1 className="text-2xl font-bold text-gray-900 mb-2">{video.title}</h1>
                  <div className="flex items-center space-x-4 text-sm text-gray-600">
                    <span>By {video.uploaded_by}</span>
                    <span>•</span>
                    <span>{new Date(video.created_at).toLocaleDateString()}</span>
                    <span>•</span>
                    <span>{formatFileSize(video.size)}</span>
                    {video.duration > 0 && (
                      <>
                        <span>•</span>
                        <span>{formatDuration(video.duration)}</span>
                      </>
                    )}
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={handleShare}
                    className="btn-secondary"
                  >
                    <ShareIcon className="h-4 w-4" />
                  </button>
                  {video.status === 'ready' && (
                    <a
                      href={`/api/v1/videos/${video.id}/stream`}
                      download
                      className="btn-secondary"
                    >
                      <ArrowDownTrayIcon className="h-4 w-4" />
                    </a>
                  )}
                </div>
              </div>

              {video.description && (
                <div className="mb-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-2">Description</h3>
                  <p className="text-gray-700 whitespace-pre-wrap">{video.description}</p>
                </div>
              )}

              {/* Video Formats */}
              {video.formats && video.formats.length > 0 && (
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-3">Available Formats</h3>
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                    {video.formats.map((format, index) => (
                      <div key={index} className="bg-gray-50 rounded-lg p-3">
                        <div className="font-medium text-gray-900">{format.quality}</div>
                        <div className="text-sm text-gray-600">{formatFileSize(format.size)}</div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Job Status Sidebar */}
          <div className="lg:min-w-80">
            <JobStatus jobs={jobs} loading={jobsLoading} />
          </div>
        </div>
      </main>
    </div>
  );
};

export default VideoDetailPage; 