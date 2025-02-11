# AdaptiveBitrateStreaming
- [ ] Create an API to upload raw video file to AWS S3 bucket
- [ ] Add task to transcode the raw file to a queue service
- [ ] Create a consumer for the queue
- [ ] Trigger a computation container (ECS container) to transcode the file using ffmpeg
- [ ] Send the transcoded segments to buckets
- [ ] Create a frontend to test the service
- [ ] Allow support for multiple cloud services
- [ ] Roadmap for predicting the bitrate

## Architecture

![image](https://github.com/user-attachments/assets/3a979972-3f6e-47f5-a04c-59c872dcd492)
