import React, { createContext, useContext, useState } from "react";

export interface Message {
  id: string;
  text: string;
  senderId: string;
  timestamp: number;
  read: boolean;
}

export interface Match {
  id: string;
  name: string;
  avatar?: string;
  messages: Message[];
}

interface MatchContextData {
  matches: Match[];
  sendMessage: (matchId: string, text: string) => void;
  markAsRead: (matchId: string) => void;
  unmatch: (matchId: string) => void;
}

const MatchContext = createContext<MatchContextData | null>(null);

export const MatchProvider = ({ children }: any) => {
  const [matches, setMatches] = useState<Match[]>([
    {
      id: "1",
      name: "Maria",
      messages: [],
    },
  ]);

  const sendMessage = (matchId: string, text: string) => {
    setMatches((prev) =>
      prev.map((match) =>
        match.id === matchId
          ? {
              ...match,
              messages: [
                ...match.messages,
                {
                  id: Date.now().toString(),
                  text,
                  senderId: "me",
                  timestamp: Date.now(),
                  read: false,
                },
              ],
            }
          : match
      )
    );
  };

  const markAsRead = (matchId: string) => {
    setMatches((prev) =>
      prev.map((match) =>
        match.id === matchId
          ? {
              ...match,
              messages: match.messages.map((msg) => ({
                ...msg,
                read: true,
              })),
            }
          : match
      )
    );
  };

  const unmatch = (matchId: string) => {
    setMatches((prev) => prev.filter((m) => m.id !== matchId));
  };

  return (
    <MatchContext.Provider
      value={{ matches, sendMessage, markAsRead, unmatch }}
    >
      {children}
    </MatchContext.Provider>
  );
};

export const useMatch = () => {
  const context = useContext(MatchContext);
  if (!context) {
    throw new Error("useMatch must be used within MatchProvider");
  }
  return context;
};