import { Ionicons } from "@expo/vector-icons";
import { useLocalSearchParams, useRouter } from "expo-router";
import { useEffect, useState } from "react";
import {
    FlatList,
    KeyboardAvoidingView,
    Platform,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from "react-native";
import { useMatch } from "../../context/MatchContext";

export default function ChatScreen() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const { matches, sendMessage, markAsRead, unmatch } = useMatch();
  const [text, setText] = useState("");

  const match = matches.find((m) => m.id === id);

  useEffect(() => {
    if (id) {
      markAsRead(id as string);
    }
  }, []);

  if (!match) return <Text>Match não encontrado</Text>;

  return (
    <KeyboardAvoidingView
      style={{ flex: 1, backgroundColor: "#f5f5f5" }}
      behavior={Platform.OS === "ios" ? "padding" : undefined}
    >
      {/* Header */}
      <View
        style={{
          flexDirection: "row",
          alignItems: "center",
          padding: 16,
          backgroundColor: "#fff",
          borderBottomWidth: 1,
          borderColor: "#ddd",
        }}
      >
        <TouchableOpacity onPress={() => router.back()}>
          <Ionicons name="arrow-back" size={24} color="black" />
        </TouchableOpacity>

        <Text
          style={{
            fontSize: 18,
            fontWeight: "bold",
            marginLeft: 16,
          }}
        >
          {match.name}
        </Text>
      </View>

      {/* Mensagens */}
      <FlatList
        contentContainerStyle={{ padding: 16 }}
        data={match.messages}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View
            style={{
              alignSelf:
                item.senderId === "me" ? "flex-end" : "flex-start",
              backgroundColor:
                item.senderId === "me" ? "#3b82f6" : "#e5e5ea",
              padding: 10,
              borderRadius: 16,
              marginVertical: 4,
              maxWidth: "75%",
            }}
          >
            <Text
              style={{
                color: item.senderId === "me" ? "#fff" : "#000",
              }}
            >
              {item.text}
            </Text>
          </View>
        )}
      />

      {/* Input */}
      <View
        style={{
          flexDirection: "row",
          padding: 12,
          backgroundColor: "#fff",
          borderTopWidth: 1,
          borderColor: "#ddd",
        }}
      >
        <TextInput
          style={{
            flex: 1,
            borderWidth: 1,
            borderColor: "#ddd",
            borderRadius: 20,
            paddingHorizontal: 16,
            paddingVertical: 8,
            marginRight: 8,
            backgroundColor: "#f9f9f9",
          }}
          value={text}
          onChangeText={setText}
          placeholder="Digite uma mensagem..."
        />

        <TouchableOpacity
          onPress={() => {
            if (text.trim()) {
              sendMessage(match.id, text);
              setText("");
            }
          }}
        >
          <Ionicons name="send" size={24} color="#3b82f6" />
        </TouchableOpacity>
      </View>

      {/* Desfazer Match */}
      <TouchableOpacity
        onPress={() => {
          unmatch(match.id);
          router.back();
        }}
        style={{
          padding: 16,
          alignItems: "center",
        }}
      >
        <Text style={{ color: "red", fontWeight: "600" }}>
          Desfazer Match
        </Text>
      </TouchableOpacity>
    </KeyboardAvoidingView>
  );
}