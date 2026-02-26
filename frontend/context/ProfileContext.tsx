import React, { createContext, ReactNode, useContext, useState } from "react";

type ProfileType = {
  name?: string;
  email?: string;
  userType?: "volunteer" | "senior";
  dateOfBirth?: string;
  gender?: string;
  phone?: string;
  bio?: string;
  address?: string;
  skills?: string[];
  needs?: string[];
  interests?: string[];
  availability?: string;
};

type ProfileContextType = {
  profile: ProfileType;
  updateProfile: (data: Partial<ProfileType>) => void;
};

const ProfileContext = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children }: { children: ReactNode }) {
  const [profile, setProfile] = useState<ProfileType>({});

  function updateProfile(data: Partial<ProfileType>) {
    setProfile((prev) => ({ ...prev, ...data }));
  }

  return (
    <ProfileContext.Provider value={{ profile, updateProfile }}>
      {children}
    </ProfileContext.Provider>
  );
}

export function useProfile() {
  const context = useContext(ProfileContext);
  if (!context) {
    throw new Error("useProfile must be used within a ProfileProvider");
  }
  return context;
}