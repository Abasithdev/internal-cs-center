import { defineStore } from "pinia";
import { Login } from "@/api/authApi";

interface AuthState{
    token: string | null;
    role: string | null;
    email: string | null;
}

export const useAuthStore = defineStore("auth",{
    state: (): AuthState => ({
        token: localStorage.getItem("token"),
        role: localStorage.getItem("role"),
        email: localStorage.getItem("email"),
    }),
    actions: {
        async login(email: string, password: string){
            const response = await Login(email, password);
            this.token = response.token;
            this.role = response.role;
            this.email = email;
            localStorage.setItem("token", response.token);
            localStorage.setItem("role", response.role);
            localStorage.setItem("email", email);
        },
        logout() {
            this.token=null;
            this.role=null;
            this.email=null;
            localStorage.clear();
        }
    }
})