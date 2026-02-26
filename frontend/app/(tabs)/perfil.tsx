import { useEffect, useState } from "react";
import {
    ScrollView,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from "react-native";
import { useProfile } from "../../context/ProfileContext";
import { BackendAdapter, UserProfile } from "@/context/adapter/BackendAdapter";
import { BackendService } from "@/context/service/BackedService";

const backendService: BackendAdapter = new BackendService();

export default function PerfilScreen() {
  // const { profile, updateProfile } = useProfile();
  const [profile, updateProfile] = useState<UserProfile | null>(null);

   useEffect(() => {
    async function loadProfile() {
        try {
          const profile = await backendService.getUserProfile();
          updateProfile(profile);

        } catch (error) {
          console.error(error);
        }
      }
    loadProfile();
  }, []);

  const  handleEdit = async () => {
    // console.log("editando", { ...profile, bio: profile?.bio, phone: profile?.phone });
    const body = {bio: profile?.bio, phone: profile?.phone}
     try {
        const profile = await backendService.updateUserProfile(body);
        console.log(profile, "editando")
        updateProfile(profile);
        alert("Perfil atualizado com sucesso!");

      } catch (error) {
        console.error(error);
        alert("Erro ao atualizar perfil!");
      }
  };

  if (!profile) {
    return <Text>Carregando...</Text>;
  }
  return (
    <ScrollView
      style={{ flex: 1, backgroundColor: "#f5f5f5" }}
      contentContainerStyle={{ padding: 24 }}
    >
      <Text style={{ fontSize: 28, fontWeight: "bold", marginBottom: 24 }}>
        Meu Perfil
      </Text>

      <Text style={{ marginBottom: 6, fontWeight: "600" }}>Email</Text>
      <Text style={{ marginBottom: 16 }}>
        {profile.email}
      </Text>


      <Text style={{ marginBottom: 6, fontWeight: "600" }}>Tipo</Text>
      <Text style={{ marginBottom: 16 }}>
        {profile.user_type === "VOLUNTEER"
          ? "Voluntário"
          : profile.user_type === "ELDERLY"
          ? "Idoso"
          : "Não definido"}
      </Text>

      {/* <Text style={{ marginBottom: 6, fontWeight: "600" }}>
        Data de Nascimento
      </Text> */}
      {/* <Text style={{ marginBottom: 16 }}>
        {profile.dateOfBirth || "Não informado"}
      </Text> */}

      {/* <Text style={{ marginBottom: 6, fontWeight: "600" }}>Endereço</Text> */}
      {/* <Text style={{ marginBottom: 16 }}>
        {profile.address || "Não informado"}
      </Text> */}

      <View style={{ marginBottom: 16 }}>
        <Text style={{ marginBottom: 6, fontWeight: "600" }}>Telefone</Text>
        <TextInput
          value={profile.phone}
          onChangeText={text =>
            updateProfile(prev => prev ? { ...prev, phone: text } : prev)
          }
          style={{
            backgroundColor: "#fff",
            padding: 12,
            borderRadius: 10,
            borderWidth: 1,
            borderColor: "#ddd",
          }}
        />
      </View>

      <View style={{ marginBottom: 24 }}>
        <Text style={{ marginBottom: 6, fontWeight: "600" }}>Bio</Text>
        <TextInput
          value={profile.bio}
          onChangeText={text =>
            updateProfile(prev => prev ? { ...prev, bio: text } : prev)
          }
          multiline
          numberOfLines={4}
          style={{
            backgroundColor: "#fff",
            padding: 12,
            borderRadius: 10,
            borderWidth: 1,
            borderColor: "#ddd",
            textAlignVertical: "top",
          }}
        />
      </View>

      <TouchableOpacity
        onPress={handleEdit}
        style={{
          backgroundColor: "#3b82f6",
          padding: 16,
          borderRadius: 12,
        }}
      >
        <Text
          style={{
            color: "#fff",
            textAlign: "center",
            fontWeight: "bold",
            fontSize: 16,
          }}
        >
          Salvar Alterações
        </Text>
      </TouchableOpacity>
    </ScrollView>
  );
}