# AdaptiveBitrateStreaming
This service allows users to **upload videos** and stream them using **Adaptive Bitrate (ABR) Streaming**.

## Architecture

![image](https://github.com/user-attachments/assets/80e0408e-3e5b-40ea-9969-5edbfc6892d7)


## 📌 Features  
- ✅ **User Authentication** (Register/Login)  
- ✅ **Secure Video Upload with Pre-signed URLs**  
- ✅ **HLS Streaming for ABR Playback**  
- ✅ **Token-based Access Control**  

---

## Design Patterns Used

### Strategy Pattern for Initializing Upload and Creating a Presign URL

I utilized the **Strategy Pattern** to handle both **Single File Upload** and **Multipart Upload** dynamically.

### Singleton Pattern

The Singleton Pattern is used in this project to ensure a single instance of critical resources such as - Database Connection, AWS Session, RabbitMQ Connection


## 🛠️ Getting Started  

### 1️⃣ Register & Login to Get Access Token  
You must **register** and **log in** to get an `access_token`.

#### 🔹 Register  
Use the following command to register a new user:  
```sh
curl -X POST "http://localhost:8080/api/v1/register" \
     -H "Content-Type: application/json" \
     -d '{"username": "your_username", "password": "your_password"}'

```


### 2 Request a presign URL
```
curl -X POST "http://localhost:8080/api/v1/initialize_upload" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer your_jwt_token_here" \
     -d '{"bucket": "abr-raw", "name": "test.mp4", "file_type": "video/mp4", "strategy": "single"}'
```
#### 3️⃣ Upload the Video Using cURL

```
curl -X PUT "https://your-storage-provider/upload-url" \
     -T "/path/to/your/video.mp4" \
     -H "Content-Type: video/mp4"
```

## TO DO

## ✅ TO-DO List  

- [ ] Create a **Dashboard** to view the status of the transcoded file  
- [ ] Implement **something similar to Netflix VMAF** for video quality assessment  





