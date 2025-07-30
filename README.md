# 🚀 Nimbus Drive - Google Drive Clone

A modern, full-stack Google Drive clone built with **Vue.js 3** and **Go (Gin)**, featuring a beautiful UI, real-time file management, and cloud storage capabilities.

![Nimbus Drive](https://img.shields.io/badge/Vue.js-3.5.17-4FC08D?style=for-the-badge&logo=vue.js)
![Go](https://img.shields.io/badge/Go-1.24.1-00ADD8?style=for-the-badge&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=for-the-badge&logo=mysql)
![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?style=for-the-badge&logo=redis)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4.17-38B2AC?style=for-the-badge&logo=tailwind-css)

## ✨ Features

### 📁 **File Management**
- **Upload/Download** files with drag & drop support
- **Create folders** and organize content hierarchically
- **Move/Copy** files and folders between locations
- **Rename** files and folders with real-time updates
- **Delete** with soft delete (trash system)
- **Bulk operations** for multiple files/folders

### ⭐ **Starred System**
- **Star/Unstar** files and folders for quick access
- **Dedicated starred view** with all starred items
- **Real-time status updates** across the application
- **Search within starred items**

### 🗑️ **Trash Management**
- **Soft delete** system with 30-day retention
- **Restore** deleted files and folders
- **Permanent delete** option
- **Bulk restore/delete** operations

### 🎨 **Modern UI/UX**
- **Responsive design** with Tailwind CSS
- **Dark/Light theme** toggle
- **Grid and List view** modes
- **Drag & drop** file uploads
- **Context menus** for quick actions
- **Breadcrumb navigation**
- **Persistent navigation** state (URL-based)

### 🔐 **Authentication & Security**
- **JWT-based authentication**
- **User registration and login**
- **Password hashing** with bcrypt
- **Protected routes** and API endpoints
- **File access control** per user

### 🔍 **Search & Organization**
- **Real-time search** across files and folders
- **File type detection** with icons
- **File size formatting**
- **Date formatting** for modifications
- **Sorting and filtering** options

## 🏗️ Architecture

### **Frontend (Vue.js 3)**
```
client/
├── src/
│   ├── components/
│   │   ├── Dashboard.vue          # Main dashboard
│   │   ├── Login.vue              # Authentication
│   │   ├── Register.vue           # User registration
│   │   └── dashboard/
│   │       ├── StarredView.vue    # Starred items view
│   │       ├── TrashView.vue      # Trash management
│   │       ├── Toolbar.vue        # Top toolbar
│   │       ├── Sidebar.vue        # Navigation sidebar
│   │       ├── ContentArea.vue    # File display
│   │       ├── UploadModal.vue    # File upload
│   │       ├── MoveModal.vue      # Move operations
│   │       └── ...                # Other UI components
│   ├── composables/
│   │   ├── useFileManager.ts      # File operations
│   │   ├── useStarred.ts          # Starred functionality
│   │   ├── useFileData.ts         # File data management
│   │   └── ...                    # Other composables
│   └── router/
│       └── index.ts               # Vue Router setup
```

### **Backend (Go/Gin)**
```
server/
├── internal/
│   ├── models/                    # Database models
│   ├── handlers/                  # API handlers
│   ├── routes/                    # Route definitions
│   ├── middleware/                # Custom middleware
│   ├── database/                  # Database setup
│   └── validation/                # Input validation
├── migrations/                    # Database migrations
├── uploads/                       # File storage
└── cmd/
    └── main.go                    # Application entry
```

## 🚀 Quick Start

### **Prerequisites**
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.24+ (for local development)

### **Using Docker (Recommended)**

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/nimbus-drive.git
   cd nimbus-drive
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Database: localhost:3307 (MySQL)
   - Redis: localhost:6379

### **Local Development**

1. **Backend Setup**
   ```bash
   cd server
   go mod download
   cp .env.example .env
   # Edit .env with your database credentials
   go run cmd/main.go
   ```

2. **Frontend Setup**
   ```bash
   cd client
   npm install
   cp .env.example .env
   # Edit .env with your API URL
   npm run dev
   ```

## 📋 API Endpoints

### **Authentication**
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### **Files**
- `GET /api/files` - List files
- `POST /api/files/upload` - Upload file
- `PUT /api/files/:id` - Update file
- `DELETE /api/files/:id` - Delete file
- `GET /api/files/trash` - List trashed files
- `POST /api/files/:id/restore` - Restore file
- `DELETE /api/files/:id/permanent` - Permanent delete

### **Folders**
- `GET /api/folders` - List folders
- `POST /api/folders` - Create folder
- `PUT /api/folders/:id` - Update folder
- `DELETE /api/folders/:id` - Delete folder
- `PATCH /api/folders/:id/move` - Move folder

### **Starred Items**
- `GET /api/starred` - List starred items
- `POST /api/starred/files/:id/star` - Star file
- `DELETE /api/starred/files/:id/star` - Unstar file
- `POST /api/starred/folders/:id/star` - Star folder
- `DELETE /api/starred/folders/:id/star` - Unstar folder

## 🛠️ Technology Stack

### **Frontend**
- **Vue.js 3** - Progressive JavaScript framework
- **Vue Router 4** - Official router for Vue.js
- **Tailwind CSS** - Utility-first CSS framework
- **Vite** - Next generation frontend tooling
- **TypeScript** - Type-safe JavaScript

### **Backend**
- **Go 1.24** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for Go
- **JWT** - JSON Web Tokens for authentication
- **Redis** - In-memory data structure store

### **Database & Storage**
- **MySQL 8** - Relational database
- **Redis** - Caching and session storage
- **File System** - Local file storage

### **DevOps**
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Air** - Live reload for Go development

## 🎯 Key Features

### **File Operations**
- ✅ Drag & drop file uploads
- ✅ Multi-file upload support
- ✅ File type detection and icons
- ✅ File size and date display
- ✅ Move/copy between folders
- ✅ Rename files and folders
- ✅ Bulk operations

### **User Experience**
- ✅ Persistent navigation state
- ✅ Real-time search
- ✅ Grid and list view modes
- ✅ Dark/light theme toggle
- ✅ Responsive design
- ✅ Context menus
- ✅ Breadcrumb navigation

### **Security**
- ✅ JWT authentication
- ✅ Password hashing
- ✅ Protected routes
- ✅ File access control
- ✅ Input validation

### **Advanced Features**
- ✅ Starred items system
- ✅ Trash management
- ✅ Soft delete with retention
- ✅ Bulk operations
- ✅ File preview (images)
- ✅ Upload progress tracking

## 🔧 Configuration

### **Environment Variables**

**Backend (.env)**
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=nimbus
DB_PASSWORD=secret
DB_NAME=nimbus_drive
REDIS_ADDR=localhost:6379
JWT_SECRET=your-super-secret-jwt-key
PORT=8080
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760
```

**Frontend (.env)**
```env
VITE_API_URL=http://localhost:8080
```

## 📁 Project Structure

```
nimbus-drive/
├── client/                 # Vue.js frontend
│   ├── src/
│   │   ├── components/    # Vue components
│   │   ├── composables/   # Vue composables
│   │   ├── router/        # Vue Router
│   │   └── assets/        # Static assets
│   ├── public/            # Public assets
│   └── package.json       # Frontend dependencies
├── server/                # Go backend
│   ├── internal/          # Application code
│   ├── migrations/        # Database migrations
│   ├── uploads/           # File storage
│   └── go.mod            # Go dependencies
├── docker-compose.yml     # Docker orchestration
└── README.md             # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Vue.js** team for the amazing framework
- **Gin** team for the fast HTTP framework
- **Tailwind CSS** for the utility-first CSS framework
- **Google Drive** for the inspiration

---

**Built with ❤️ using Vue.js and Go** 