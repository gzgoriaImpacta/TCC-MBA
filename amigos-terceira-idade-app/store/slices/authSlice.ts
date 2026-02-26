import { createSlice, PayloadAction } from '@reduxjs/toolkit';

// 1. Tipagem do Estado
interface AuthState {
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean; // Útil para mostrar um "loading" enquanto verificamos o token ao abrir o app
}

// 2. Estado Inicial
const initialState: AuthState = {
  token: null,
  isAuthenticated: false,
  isLoading: true,
};

// 3. Criação do Slice
export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    // Ação para quando o login tem sucesso ou quando carregamos o token salvo
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload;
      state.isAuthenticated = true;
      state.isLoading = false;
    },
    // Ação para logout
    clearToken: (state) => {
      state.token = null;
      state.isAuthenticated = false;
      state.isLoading = false;
    },
    // Ação para finalizar o carregamento inicial (caso não haja token)
    setLoadingComplete: (state) => {
      state.isLoading = false;
    },
  },
});

export const { setToken, clearToken, setLoadingComplete } = authSlice.actions;
export default authSlice.reducer;