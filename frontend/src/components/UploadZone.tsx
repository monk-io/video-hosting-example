import React, { useState, useRef, DragEvent } from 'react';
import { CloudArrowUpIcon, VideoCameraIcon } from '@heroicons/react/24/outline';
import { useUpload } from '../hooks/useUpload';
import ProgressBar from './ProgressBar';

interface UploadZoneProps {
  onUploadComplete?: (videoId: string) => void;
}

const UploadZone: React.FC<UploadZoneProps> = ({ onUploadComplete }) => {
  const [isDragOver, setIsDragOver] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const {
    uploadState,
    setFile,
    setTitle,
    setDescription,
    uploadVideo,
    resetUpload,
    validateFile,
  } = useUpload();

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOver(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileSelect = (file: File) => {
    const error = validateFile(file);
    if (error) {
      alert(error);
      return;
    }
    setFile(file);
  };

  const handleUpload = async () => {
    const videoId = await uploadVideo();
    if (videoId && onUploadComplete) {
      onUploadComplete(videoId);
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  if (uploadState.success) {
    return (
      <div className="card p-6 text-center">
        <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
          <VideoCameraIcon className="h-6 w-6 text-green-600" />
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">Upload Successful!</h3>
        <p className="text-sm text-gray-500 mb-4">
          Your video has been uploaded and is being processed.
        </p>
        <button
          onClick={resetUpload}
          className="btn-secondary"
        >
          Upload Another Video
        </button>
      </div>
    );
  }

  return (
    <div className="card p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Upload Video</h2>
      
      {!uploadState.file ? (
        <div
          className={`upload-zone ${isDragOver ? 'drag-over' : ''}`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="video/*"
            onChange={handleFileInputChange}
            className="hidden"
          />
          <CloudArrowUpIcon className="mx-auto h-12 w-12 text-gray-400 mb-4" />
          <p className="text-lg font-medium text-gray-900 mb-2">
            Drag and drop your video here
          </p>
          <p className="text-sm text-gray-500 mb-4">
            or click to select a file
          </p>
          <p className="text-xs text-gray-400">
            Supports MP4, AVI, MOV, WMV, WebM (max 1GB)
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {/* File Info */}
          <div className="flex items-center p-4 bg-gray-50 rounded-lg">
            <VideoCameraIcon className="h-8 w-8 text-gray-400 mr-3" />
            <div className="flex-1">
              <p className="text-sm font-medium text-gray-900">
                {uploadState.file.name}
              </p>
              <p className="text-xs text-gray-500">
                {formatFileSize(uploadState.file.size)}
              </p>
            </div>
            <button
              onClick={() => setFile(null)}
              className="text-gray-400 hover:text-gray-600"
              disabled={uploadState.isUploading}
            >
              ×
            </button>
          </div>

          {/* Video Details Form */}
          <div className="space-y-4">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                Title *
              </label>
              <input
                type="text"
                id="title"
                value={uploadState.title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter video title"
                className="w-full rounded-md border-gray-300 shadow-sm focus:border-red-500 focus:ring-red-500"
                disabled={uploadState.isUploading}
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                id="description"
                rows={3}
                value={uploadState.description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Enter video description (optional)"
                className="w-full rounded-md border-gray-300 shadow-sm focus:border-red-500 focus:ring-red-500"
                disabled={uploadState.isUploading}
              />
            </div>
          </div>

          {/* Upload Progress */}
          {uploadState.isUploading && (
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-700">Uploading...</span>
                <span className="text-gray-500">{uploadState.progress.percentage}%</span>
              </div>
              <ProgressBar 
                progress={uploadState.progress.percentage} 
                color="bg-red-500"
              />
              <p className="text-xs text-gray-500">
                {formatFileSize(uploadState.progress.loaded)} of {formatFileSize(uploadState.progress.total)}
              </p>
            </div>
          )}

          {/* Error Message */}
          {uploadState.error && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-md">
              <p className="text-sm text-red-700">{uploadState.error}</p>
            </div>
          )}

          {/* Upload Button */}
          <button
            onClick={handleUpload}
            disabled={uploadState.isUploading || !uploadState.title.trim()}
            className="btn-primary w-full"
          >
            {uploadState.isUploading ? 'Uploading...' : 'Upload Video'}
          </button>
        </div>
      )}
    </div>
  );
};

export default UploadZone; 