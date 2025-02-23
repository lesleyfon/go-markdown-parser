## Markdown Note-taking App

Build a simple note-taking app that lets users upload markdown files, check the grammar, save the note, and render it in HTML. 
The goal of this project is to help you learn how to handle file uploads in a RESTful API, parse and render markdown files using libraries, and check the grammar of the notes.

[https://roadmap.sh/projects/markdown-note-taking-app](https://roadmap.sh/projects/markdown-note-taking-app)

![alt text](image.png)


### Requirements

You have to implement the following features:
- You'll provide an endpoint to check the grammar of the note.
- You'll also provide an endpoint to save the note that can be passed in as Markdown text.
- Provide an endpoint to list the saved notes (i.e. uploaded markdown files).
- Return the HTML version of the Markdown note (rendered note) through another endpoint.

## Architecture

### Backend (Go)
- **Framework**: Gin Web Framework
- **Database**: MongoDB
- **Authentication**: JWT-based authentication
- **Key Libraries**:
  - `github.com/sajari/fuzzy` for spell checking
  - `github.com/yuin/goldmark` for markdown to HTML conversion
  - `go.mongodb.org/mongo-driver` for MongoDB operations

### Frontend (React + TypeScript)
- **Framework**: React with TypeScript
- **Router**: TanStack Router
- **State Management**: TanStack Query
- **UI Components**: Custom components with Tailwind CSS

## Features

### 1. Authentication
- User signup and login functionality
- JWT-based authentication with refresh tokens
- Protected routes and API endpoints

### 2. Spell Checking
- Real-time spell checking of markdown files
- Fuzzy matching algorithm for suggestions
- Parallel processing for improved performance
- Levenshtein distance calculation for accurate suggestions
- Custom dictionary support

### 3. File Management
- Upload markdown files
- Save files to MongoDB
- List all files for authenticated users
- View individual file details
- File versioning with timestamps

### 4. Markdown Processing
- Convert markdown to HTML
- Highlight misspelled words
- Interactive UI for viewing suggestions
- Support for code blocks and other markdown features

## API Endpoints

### Authentication
```
POST /auth/v1/signup - Register new user
POST /auth/v1/login - User login
GET /auth/v1/authenticate - Verify authentication
```
### Markdown Operations
```
POST /api/v1/markdown - Upload and spell check markdown file
GET /api/v1/markdown/files - Get all files for authenticated user
GET /api/v1/markdown/files/:file_id - Get specific file by ID
```


## Technical Details

### Spell Checking Algorithm
- Uses fuzzy matching with configurable threshold and depth
- Parallel processing of words in chunks
- Custom dictionary support with case-insensitive matching
- Levenshtein distance filtering for accurate suggestions

```go
type SpellCheckConfig struct {
    LevenshteinThreshold int
    FuzzyModelDepth      int
    FuzzyModelThreshold  int
}
```

### File Processing
- Maximum file size: 8MB
- Supported format: Markdown (.md)
- Automatic HTML conversion
- Real-time spell checking
- Base64 encoding for HTML content transfer

### Security Features
- JWT-based authentication
- Token refresh mechanism
- Password hashing
- CORS configuration
- Request timeout handling

## Frontend Components

### Main Components
- `AppSideBar`: Navigation and file listing
- `FileDetails`: Individual file viewer
- `Index`: Main file upload and spell check interface

### State Management
- React Query for server state
- Local state for UI components
- Memoized components for performance

## Configuration

### Environment Variables

```bash
PORT=8080 (default)
SECRET_KEY=your_jwt_secret
MONGODB_URL=your_mongodb_url
```


### CORS Configuration
```go
AllowOrigins: ["http://localhost:5173"]
AllowMethods: ["PUT", "PATCH", "POST", "GET", "OPTIONS"]
AllowHeaders: ["Origin", "Authorization", "Content-Type", "Accept"]
```

## Development Setup

### Prerequisites
- Go 1.23.2 or higher
- Node.js and npm
- MongoDB instance
- Git

### Installation
1. Clone the repository
2. Set up environment variables
3. Install backend dependencies: `go mod download`
4. Install frontend dependencies: `cd client && npm install`
5. Start backend server: `go run main.go`
6. Start frontend development server: `cd client && npm run dev`

## Todo
- [ ] For file saving
  - [ ] Add file size limits
  - [ ] Add content-type validation
  - [ ] Add user quotas (max files per user)
  - [ ] Add indexes on file_name and user_id fields in MongoDB
  - [ ] Add file versioning if needed
  - [ ] Add request timeout handling
- [ ] Stream response instead of returning all at once
- [ ] Improve spell checking performance
- [ ] Add support for custom dictionaries

## License


**Todo:**
- [X] Create a simple gin server
- [X] Add a simple home endpoint that response with a message `Welcome to Markdown parser`
- [X] Create an endpoint that lets users upload a markdown file.
- [X] Parse the markdown file and find spelling mistakes 
- [X] Add a simple UI to upload a markdown file and see the result.
- [X] Add a simple UI to list all the uploaded markdown files.
- [X] Add a simple UI to see the HTML version of the Markdown note.
- [ ] For file saving
  - [ ] Add file size limits
  - [ ] Add content-type validation
  - [ ] Add user quotas (max files per user)
  - [ ] Add indexes on file_name and user_id fields in MongoDB
  - [ ] Add file versioning if needed
  - [ ] Add request timeout handling
