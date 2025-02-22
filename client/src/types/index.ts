// Auth types
export interface AuthUser {
    email: string;
    refresh_token: string;
    token: string;
    user_id: string;
}

export interface AuthResponse {
    isAuthenticated: boolean;
    message: string;
    status: number;
    user: AuthUser;
}



export interface UnAuthenticateUserResponse {
	message: string;
	status: number;
	isAuthenticated: boolean;
}
// File types
export interface FileType {
    _id: string;
    file_id: string;
    created_at: string;
    file_content: string;
    file_name: string;
    updated_at: string;
    user_id: string;
}

export interface FilesResponse {
    files: FileType[];
    message: string;
    status: number;
} 
