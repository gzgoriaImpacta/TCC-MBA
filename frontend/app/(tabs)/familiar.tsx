// import { useState } from "react";
// import {
//     FlatList,
//     StyleSheet,
//     Text,
//     TouchableOpacity,
//     View,
// } from "react-native";

// type Mensagem = {
//   id: string;
//   nome: string;
//   mensagem: string;
//   lida: boolean;
// };

// export default function Familiar() {
//   const [mensagens, setMensagens] = useState<Mensagem[]>([
//     {
//       id: "1",
//       nome: "Maria (Filha)",
//       mensagem: "Oi mãe! Passando para dizer que te amo ❤️",
//       lida: false,
//     },
//     {
//       id: "2",
//       nome: "João (Filho)",
//       mensagem: "Estarei aí no domingo!",
//       lida: false,
//     },
//     {
//       id: "3",
//       nome: "Neto Pedro",
//       mensagem: "Vó, fiz um desenho pra você 😊",
//       lida: true,
//     },
//   ]);

//   const marcarComoLida = (id: string) => {
//     setMensagens((prev) =>
//       prev.map((msg) =>
//         msg.id === id ? { ...msg, lida: true } : msg
//       )
//     );
//   };

//   const naoLidas = mensagens.filter((m) => !m.lida).length;

//   return (
//     <View style={styles.container}>
//       <Text style={styles.title}>Mensagens da Família</Text>

//       {naoLidas > 0 && (
//         <View style={styles.badge}>
//           <Text style={styles.badgeText}>
//             {naoLidas} nova(s) mensagem(ns)
//           </Text>
//         </View>
//       )}

//       <FlatList
//         data={mensagens}
//         keyExtractor={(item) => item.id}
//         renderItem={({ item }) => (
//           <TouchableOpacity
//             style={[
//               styles.card,
//               { backgroundColor: item.lida ? "#EAEAEA" : "#fff" },
//             ]}
//             onPress={() => marcarComoLida(item.id)}
//           >
//             <Text style={styles.nome}>{item.nome}</Text>
//             <Text style={styles.mensagem}>{item.mensagem}</Text>

//             {!item.lida && (
//               <View style={styles.naoLido}>
//                 <Text style={styles.naoLidoTexto}>Não lida</Text>
//               </View>
//             )}
//           </TouchableOpacity>
//         )}
//       />
//     </View>
//   );
// }

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     padding: 20,
//     backgroundColor: "#F4F8FB",
//   },
//   title: {
//     fontSize: 26,
//     fontWeight: "bold",
//     marginBottom: 15,
//   },
//   badge: {
//     backgroundColor: "#FF4D4D",
//     padding: 8,
//     borderRadius: 10,
//     marginBottom: 15,
//     alignSelf: "flex-start",
//   },
//   badgeText: {
//     color: "#fff",
//     fontWeight: "bold",
//   },
//   card: {
//     padding: 16,
//     borderRadius: 14,
//     marginBottom: 12,
//     elevation: 2,
//   },
//   nome: {
//     fontSize: 18,
//     fontWeight: "bold",
//   },
//   mensagem: {
//     fontSize: 16,
//     marginTop: 4,
//   },
//   naoLido: {
//     marginTop: 8,
//     alignSelf: "flex-start",
//     backgroundColor: "#2E86DE",
//     paddingHorizontal: 8,
//     paddingVertical: 4,
//     borderRadius: 8,
//   },
//   naoLidoTexto: {
//     color: "#fff",
//     fontSize: 12,
//   },
// });