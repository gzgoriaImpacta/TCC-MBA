import { useRouter } from "expo-router";
import { useState } from "react";
import { Text, View } from "react-native";
import { useProfile } from "../context/ProfileContext";

export default function OnboardingScreen() {
  const router = useRouter();
  const { updateProfile } = useProfile();

  const [dateOfBirth, setDateOfBirth] = useState("");
  const [gender, setGender] = useState("");
  const [phone, setPhone] = useState("");
  const [bio, setBio] = useState("");

  const handleFinish = () => {
    updateProfile({
      dateOfBirth,
      gender,
      phone,
      bio,
    });

    router.replace("/(tabs)");
  };

  return (
    <View>
      <Text>Onboarding</Text>
    </View>
  );
}