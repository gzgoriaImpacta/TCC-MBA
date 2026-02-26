import { UserModel } from "@/model/UserModel";
import { api } from "./api";
import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

import {
    Appointment,
    AuthResponse,
    BackendAdapter,
    Connection,
    CreateAppointmentRequest,
    HealthCheckResponse,
    Interest,
    Invitation,
    LoginRequest,
    UpdateProfileRequest,
    UserProfile,
    UserSuggestion
} from "../adapter/BackendAdapter";

export class BackendService implements BackendAdapter {
    private baseUrl: string;
    private accessToken: string | null = null;

    constructor(baseUrl: string = "http://localhost:8080/api/v1") {
        this.baseUrl = baseUrl;
    }

    private async request<T>(method: string, endpoint: string, data?: unknown): Promise<T> {
        const headers: HeadersInit = {
            "Content-Type": "application/json",
        };

        if (this.accessToken) {
            headers["Authorization"] = `Bearer ${this.accessToken}`;
        }

        const response = await fetch(`${this.baseUrl}${endpoint}`, {
            method,
            headers,
            body: data ? JSON.stringify(data) : undefined,
        });

        if (!response.ok) {
            throw new Error(`API Error: ${response.status} ${response.statusText}`);
        }

        return response.json();
    }

    // Health Check
    async healthCheck(): Promise<HealthCheckResponse> {
        const url = `${this.baseUrl}/health`;
        console.log("Checking backend health at:", url);
        try {
            const response = await fetch(url);

            if (!response.ok) {
                throw new Error(`Health check failed: ${response.status} ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error("Health check failed:", error);
            throw error;
        }

    }

    // Registration & Authentication
    async register(data: UserModel): Promise<AuthResponse> {
        const url = `${this.baseUrl}/auth/register`;
        const body = JSON.stringify(data).replace(/userType/g, "user_type"); // Adjusting to backend's expected format
        console.log("User data being sent for registration:", body);
        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: body,
            });
            if (!response.ok) {
                throw new Error(`Registration failed: ${response.status} ${response.statusText}`);
            }
            const result = await response.json();
            this.accessToken = result.data.access_token;
            return result;
        } catch (error) {
            console.error("Registration error:", error);
            throw error;
        }
    }

    // async login(data: LoginRequest): Promise<AuthResponse> {
    //     const response = await this.request<AuthResponse>("POST", "/auth/login", data);
    //     this.accessToken = response.data.access_token;
    //     return response;
    // }
    async login(data: LoginRequest): Promise<AuthResponse> {
        const response = await api.post<AuthResponse>(`${this.baseUrl}/auth/login`, data);

        const token = response.data.data.access_token;

        if (Platform.OS === "web") {
            localStorage.setItem("user_jwt", token);
        } else {
            await SecureStore.setItemAsync("user_jwt", token);
        }
        console.log(response.data);
        return response.data;
    }

    async logout(): Promise<void> {
        if (Platform.OS === "web") {
            localStorage.removeItem("user_jwt");
        } else {
            await SecureStore.deleteItemAsync("user_jwt");
        }

        this.accessToken = null;
    }

    // Users
    async getUserProfile(): Promise<UserProfile> {
        const response = await api.get(`${this.baseUrl}/users/me`);
        console.log(response, "joaninha")
        return response.data.data;
    }

    async updateUserProfile(data: UpdateProfileRequest): Promise<UserProfile> {
        const response = await api.put(`${this.baseUrl}/users/me`, data);
        return response.data.data;
    }

    async deactivateAccount(): Promise<void> {
        await this.request("POST", "/users/deactivate");
    }

    async refreshTokens(refreshToken: string): Promise<AuthResponse> {
        const response = await this.request<AuthResponse>("POST", "/auth/refresh", { refreshToken });
        this.accessToken = response.data.access_token;
        return response;
    }

    // Appointments
    async createAppointment(data: CreateAppointmentRequest): Promise<Appointment> {
        return this.request<Appointment>("POST", "/appointments", data);
    }

    async getAppointments(): Promise<Appointment[]> {
        return this.request<Appointment[]>("GET", "/appointments");
    }

    async getUpcomingAppointments(): Promise<Appointment[]> {
        return this.request<Appointment[]>("GET", "/appointments/upcoming");
    }

    async getAppointmentDetails(appointmentId: number): Promise<Appointment> {
        return this.request<Appointment>("GET", `/appointments/${appointmentId}`);
    }

    async acceptAppointment(appointmentId: number): Promise<void> {
        await this.request("POST", `/appointments/${appointmentId}/accept`);
    }

    async declineAppointment(appointmentId: number): Promise<void> {
        await this.request("POST", `/appointments/${appointmentId}/decline`);
    }

    async cancelAppointment(appointmentId: number): Promise<void> {
        await this.request("POST", `/appointments/${appointmentId}/cancel`);
    }

    // Invitations
    async getSentInvitations(): Promise<Invitation[]> {
        return this.request<Invitation[]>("GET", "/invitations/sent");
    }

    async getReceivedInvitations(): Promise<Invitation[]> {
        return this.request<Invitation[]>("GET", "/invitations/received");
    }

    // Interests
    async getInterests(): Promise<Interest[]> {
        return this.request<Interest[]>("GET", "/interests");
    }

    async getInterestDetails(interestId: number): Promise<Interest> {
        return this.request<Interest>("GET", `/interests/${interestId}`);
    }

    // Connections/Matching
    async createConnection(userId: number): Promise<Connection> {
        return this.request<Connection>("POST", "/connections", { userId });
    }

    async getConnections(): Promise<Connection[]> {
        return this.request<Connection[]>("GET", "/connections");
    }

    async acceptConnection(connectionId: number): Promise<void> {
        await this.request("POST", `/connections/${connectionId}/accept`);
    }

    async rejectConnection(connectionId: number): Promise<void> {
        await this.request("POST", `/connections/${connectionId}/reject`);
    }

    async getSuggestions(): Promise<UserSuggestion[]> {
        return this.request<UserSuggestion[]>("GET", "/suggestions");
    }


}