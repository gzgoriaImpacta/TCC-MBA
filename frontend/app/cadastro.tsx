import type { UserModel } from "@/model/UserModel";
import { useRouter } from "expo-router";
import { useState } from "react";
import { USER_TYPES } from "../constants/userTypesConst";
import DateTimePicker from '@react-native-community/datetimepicker';
import { BackendAdapter, HealthCheckResponse } from "../context/adapter/BackendAdapter";
import { BackendService } from "../context/service/BackedService";

import {
  Button,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View
} from "react-native";
type UserType = typeof USER_TYPES[keyof typeof USER_TYPES];

export default function Cadastro() {
  const router = useRouter();
  const user = {} as UserModel;
  const backend: BackendAdapter = new BackendService();

  const [selectedDate, setSelectedDate] = useState<Date>();
  const [userType, setUserType] = useState<UserType>();
  const [step, setStep] = useState(1);

  // Básico
  const [nome, setNome] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [dateOfBirth, setDateOfBirth] = useState("");

  // Perfil
  const [bio, setBio] = useState("");
  const [interesses, setInteresses] = useState<string[]>([]);

  const [show, setShow] = useState(false);

  const showMode = () => {
    setShow(true);
  };

  const interestOptions = [
    "Leitura",
    "Caminhadas",
    "Jogos",
    "Música",
    "Companhia",
  ];

  const toggleInterest = (interest: string) => {
    setInteresses((prev) =>
      prev.includes(interest)
        ? prev.filter((i) => i !== interest)
        : [...prev, interest]
    );
  };

  const handleNext = () => {
    if (step < 3) setStep(step + 1);
  };

  const handleBack = () => {
    if (step > 1) setStep(step - 1);
    else router.back();
  };
  const handleCancel = () => {
    router.replace("/");
  };


  async function handleFinish() {
    const currentYearDate = new Date().getFullYear();
    const birthYear = Platform.OS === "web" ? new Date(dateOfBirth).getFullYear() : selectedDate?.getFullYear() || currentYearDate;
    user.name = nome;
    user.email = email;
    user.age = currentYearDate - birthYear;
    user.password = password;
    user.bio = bio;
    user.userType = userType || USER_TYPES.VOLUNTARIO;   

    const healthCheckResponse: HealthCheckResponse = await backend.healthCheck();
    
    console.log("Usuário registrado:", user);

    if (healthCheckResponse.status === "ok") {
      await backend.register(user);
      console.log("Usuário registrado:", user);
      alert("Cadastro realizado com sucesso! " + user.name);
    }
    // Aqui depois você pode salvar globalmente
    router.replace("/");
  };

  return (
    <View style={{ flex: 1, backgroundColor: "#F4F8FB" }}>
      <ScrollView contentContainerStyle={styles.container}>
        {/* Indicador */}
        <View style={styles.progressContainer}>
          {[1, 2, 3].map((s) => (
            <View
              key={s}
              style={[
                styles.progressBar,
                { backgroundColor: s <= step ? "#4CAF50" : "#ddd" },
              ]}
            />
          ))}
        </View>

        {/* STEP 1 */}
        {step === 1 && (
          <>
            <Text style={styles.title}>Cadastro</Text>

            <TextInput
              placeholder="Nome completo"
              value={nome}
              onChangeText={setNome}
              style={styles.input}
            />

            <TextInput
              placeholder="E-mail"
              value={email}
              onChangeText={setEmail}
              style={styles.input}
            />

            <TextInput
              placeholder="Senha"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              style={styles.input}
            />

            {Platform.OS === "web" && (
              <TextInput
                placeholder="Data de nascimento (DD/MM/AAAA)"
                defaultValue={new Date().toLocaleDateString()}
                onChangeText={setDateOfBirth}
                
                style={styles.input}
              />
            )}

            {Platform.OS === "android" &&  (
              <View>
                <Button onPress={showMode} title="Show date picker!" />
                <Text>selected: {selectedDate?.toLocaleString()}</Text>
                {show && (
                  <DateTimePicker
                    testID="dateTimePicker"
                    value={new Date()}
                    mode="date"
                    onChange={(event: any, date?: Date) => {
                      setSelectedDate(date);
                    }}
                  />
                )}
              </View>
            )}      
          </>
        )}

        {/* STEP 2 */}
        {step === 2 && (
          <>
            <Text style={styles.title}>Perfil</Text>

            <TextInput
              placeholder="Bio"
              value={bio}
              onChangeText={setBio}
              multiline
              style={[styles.input, { height: 100 }]}
            />

            <View style={styles.tagsContainer}>
              {interestOptions.map((interest) => (
                <TouchableOpacity
                  key={interest}
                  onPress={() => toggleInterest(interest)}
                  style={[
                    styles.tag,
                    {
                      backgroundColor: interesses.includes(
                        interest
                      )
                        ? "#4CAF50"
                        : "#fff",
                    },
                  ]}
                >
                  <Text
                    style={{
                      color: interesses.includes(interest)
                        ? "#fff"
                        : "#000",
                    }}
                  >
                    {interest}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
            
            <View>
              <Text style={styles.title}>Tipo de perfil</Text>
              {Object.entries(USER_TYPES).map(([key, value]) => (
                <TouchableOpacity
                  key={key}
                  onPress={() => setUserType(value)}
                  style={[
                    styles.primaryButton,
                    {
                      backgroundColor:
                        userType === value ? "#4CAF50" : "#fff",
                      borderColor: userType === value ? "#4CAF50" : "#ddd",
                    },
                  ]}
                >
                  <Text
                    style={{
                      color: userType === value ? "#fff" : "#000",
                      fontWeight: "600",
                    }}
                  >
                    {key.toLowerCase()}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </>
        )}

        {/* Botões */}
        <View style={styles.navButtons}>
          <TouchableOpacity
            style={styles.secondaryButton}
            onPress={handleBack}
          >
            <Text>Voltar</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={styles.secondaryButton}
            onPress={handleCancel}
          >
            <Text>Cancelar</Text>
          </TouchableOpacity>

          {step < 2 ? (
            <TouchableOpacity
              style={styles.primaryButton}
              onPress={handleNext}
            >
              <Text style={styles.buttonText}>Próximo</Text>
            </TouchableOpacity>
          ) : (
            <TouchableOpacity
              style={styles.primaryButton}
              onPress={handleFinish}
            >
              <Text style={styles.buttonText}>Finalizar</Text>
            </TouchableOpacity>
          )}
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: 24 },
  title: {
    fontSize: 26,
    fontWeight: "bold",
    marginBottom: 20,
  },
  input: {
    backgroundColor: "#fff",
    padding: 14,
    borderRadius: 12,
    marginBottom: 16,
  },
  primaryButton: {
    backgroundColor: "#4CAF50",
    padding: 16,
    borderRadius: 12,
    alignItems: "center",
    marginBottom: 10,
  },
  secondaryButton: {
    backgroundColor: "#ddd",
    padding: 16,
    borderRadius: 12,
    alignItems: "center",
  },
  buttonText: {
    color: "#fff",
    fontWeight: "bold",
  },
  success: {
    marginTop: 10,
    color: "green",
  },
  navButtons: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginTop: 20,
    gap: 10,
  },
  progressContainer: {
    flexDirection: "row",
    justifyContent: "center",
    marginBottom: 30,
  },
  progressBar: {
    width: 40,
    height: 4,
    marginHorizontal: 4,
    borderRadius: 2,
  },
  tagsContainer: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 10,
  },
  tag: {
    padding: 10,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: "#ddd",
  },
});