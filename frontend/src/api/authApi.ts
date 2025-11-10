import api from "./axiosClient";

export interface LoginResponse{
    token: string;
    role: "cs" | "operation";
}

export async function Login(email:string, password: string): Promise<LoginResponse> {
    const {data} = await api.post("/auth/login", {email, password})
    return data
}