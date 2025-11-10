export interface User{
    email: string;
    role: "cs" | "operational";
    token?: string;
}

export interface LoginRequest{
    email: string;
    password: string;
}

export interface LoginRespose{
    token: string;
    role: "cs" | "operational"
}