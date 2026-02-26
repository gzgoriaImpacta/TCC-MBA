import { UserModel } from "@/model/UserModel";

export interface BackendAdapter {
    // Health Check
    healthCheck(): Promise<HealthCheckResponse>;

    // Authentication
    register(data: UserModel): Promise<AuthResponse>;
    login(data: LoginRequest): Promise<AuthResponse>;
    refreshTokens(refreshToken: string): Promise<AuthResponse>;
    logout(): Promise<void>;

    // Appointments
    createAppointment(data: CreateAppointmentRequest): Promise<Appointment>;
    getAppointments(): Promise<Appointment[]>;
    getUpcomingAppointments(): Promise<Appointment[]>;
    getAppointmentDetails(appointmentId: number): Promise<Appointment>;
    acceptAppointment(appointmentId: number): Promise<void>;
    declineAppointment(appointmentId: number): Promise<void>;
    cancelAppointment(appointmentId: number): Promise<void>;

    // Invitations
    getSentInvitations(): Promise<Invitation[]>;
    getReceivedInvitations(): Promise<Invitation[]>;

    // Interests
    getInterests(): Promise<Interest[]>;
    getInterestDetails(interestId: number): Promise<Interest>;

    // Connections/Matching
    createConnection(userId: number): Promise<Connection>;
    getConnections(): Promise<Connection[]>;
    acceptConnection(connectionId: number): Promise<void>;
    rejectConnection(connectionId: number): Promise<void>;
    getSuggestions(): Promise<UserSuggestion[]>;

    // Users
    getUserProfile(): Promise<UserProfile>;
    updateUserProfile(data: UpdateProfileRequest): Promise<UserProfile>;
    deactivateAccount(): Promise<void>;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
    user_type: "VOLUNTEER" | "ELDERLY";
}

export interface LoginRequest {
    name?: string;
    email: string;
    password: string;
}

export interface AuthResponse {
    data: {
        access_token: string;
        refresh_token: string;
        user: UserProfile;
    };
}

export interface CreateAppointmentRequest {
    title: string;
    date: string;
    with_user_id: number;
}

export interface Appointment {
    id: number;
    title: string;
    date: string;
    with_user_id: number;
    status: string;
}

export interface Invitation {
    id: number;
    user_id: number;
    status: string;
}

export interface Interest {
    id: number;
    name: string;
}

export interface Connection {
    id: number;
    user_id: number;
    status: string;
}
export interface HealthCheckResponse {
    service: string;
    status: string;
}

export interface UserSuggestion {
    id: number;
    name: string;
    interests: Interest[];
}

export interface UserProfile {
    id: string;
    name: string;
    email: string;
    bio?: string;
    user_type: string;
    phone?: string;
}

export interface UpdateProfileRequest {
    phone?: string;
    bio?: string;
}