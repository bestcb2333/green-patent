export interface User {
  id: number,
  createdAt: Date,
  name: string,
  nickname: string,
  email: string,
  phone: string,
  admin: boolean,
  patents: Patent[],
  reports: Report[],
  keywords: Keyword[],
  suggestions: Suggestion[],
}

export interface Patent {
  id: number,
  createdAt: Date,
  name: string,
  number: string,
  status: boolean,
  file: string,
  user: User|null,
  keywords: Keyword[],
}

export interface Report {
  id: number,
  createdAt: Date,
  name: string,
  year: number,
  status: boolean,
  file: string,
  user: User|null,
  keywords: Keyword[],
}

export interface Suggestion {
  id: number,
  createdAt: Date,
  user: User|null,
  patent: Patent|null,
  report: Report|null,
  title: string,
  content: string,
}

export interface Keyword {
  id: number,
  createdAt: Date,
  user: User|null,
  value: string,
  patents: Patent[],
  reports: Report[],
}
