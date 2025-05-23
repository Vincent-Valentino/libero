// --- Interfaces (Data Structures) ---

interface Team {
  id: string | number;
  name: string;
  logo: string; // URL or path
}

interface Match {
  id: string;
  homeTeam: Team;
  awayTeam: Team;
  date: string; // ISO string or Date object
  time: string; // e.g., "15:00"
  venue: string;
  status: 'scheduled' | 'live' | 'finished';
  score?: { home: number; away: number };
}

interface Player {
  id: string;
  name: string;
  photo: string; // URL or path
  team: Team;
  position: string;
}

interface PlayerStat extends Player {
  value: number; // e.g., goals, assists, clean sheets
}

interface LeagueTableRow {
  position: number;
  team: Team;
  played: number;
  won: number;
  drawn: number;
  lost: number;
  goalsFor: number;
  goalsAgainst: number;
  goalDifference: number;
  points: number;
}

interface LeagueHistoryItem {
  season: string; // e.g., "2023/24"
  winner?: Team;
  topScorer?: PlayerStat;
  bestPlayer?: Player;
  teamOfTheSeason?: Player[];
  // Add other awards as needed
}

interface LeagueData {
  name: string;
  logo: string; // URL or path
  themeColor: string; // Hex code e.g., '#ff0000'
  upcomingMatches: Match[];
  topScorers: PlayerStat[];
  topAssists: PlayerStat[];
  mostCleanSheets: PlayerStat[]; // Assuming for goalkeepers/teams
  table: LeagueTableRow[];
  history: LeagueHistoryItem[];
}

interface LeagueMetadata {
  id: string;
  name: string;
  code: string;
  logo: string;
  themeColor: string;
}

// --- Mock Teams (Add more as needed) ---
const premierLeagueTeams: Record<string, Team> = {
  mci: { id: 'mci', name: 'Man City', logo: '/path/to/mci.png' }, // Placeholder logo
  liv: { id: 'liv', name: 'Liverpool', logo: '/path/to/liv.png' }, // Placeholder logo
  ars: { id: 'ars', name: 'Arsenal', logo: '/path/to/ars.png' }, // Placeholder logo
};
const laLigaTeams: Record<string, Team> = {
  rma: { id: 'rma', name: 'Real Madrid', logo: '/path/to/rma.png' },
  bar: { id: 'bar', name: 'Barcelona', logo: '/path/to/bar.png' },
};
const serieATeams: Record<string, Team> = {
  int: { id: 'int', name: 'Inter Milan', logo: '/path/to/int.png' },
  juv: { id: 'juv', name: 'Juventus', logo: '/path/to/juv.png' },
};
const bundesligaTeams: Record<string, Team> = {
  bay: { id: 'bay', name: 'Bayern Munich', logo: '/path/to/bay.png' },
  bvb: { id: 'bvb', name: 'Dortmund', logo: '/path/to/bvb.png' },
};
const ligue1Teams: Record<string, Team> = {
  psg: { id: 'psg', name: 'Paris SG', logo: '/path/to/psg.png' },
  mon: { id: 'mon', name: 'Monaco', logo: '/path/to/mon.png' },
};

// --- Mock League Data ---

const premierLeagueMockData: LeagueData = {
  name: 'Premier League',
  logo: '/public/Premier League.svg',
  themeColor: '#3d195b', // PL purple
  upcomingMatches: [
    { id: 'pl_m1', homeTeam: premierLeagueTeams.mci, awayTeam: premierLeagueTeams.liv, date: '2025-04-05', time: '15:00', venue: 'Etihad Stadium', status: 'scheduled' },
    { id: 'pl_m2', homeTeam: premierLeagueTeams.ars, awayTeam: premierLeagueTeams.mci, date: '2025-04-12', time: '17:30', venue: 'Emirates Stadium', status: 'scheduled' },
  ],
  topScorers: [
    { id: 'p1', name: 'Erling Haaland', photo: '/public/Erling Haaland.png', team: premierLeagueTeams.mci, position: 'FW', value: 25 },
    { id: 'p2', name: 'Mohamed Salah', photo: '/public/Mohamed Salah.png', team: premierLeagueTeams.liv, position: 'FW', value: 22 },
  ],
  topAssists: [
     { id: 'p3', name: 'Kevin De Bruyne', photo: '/path/to/kdb.png', team: premierLeagueTeams.mci, position: 'MF', value: 15 },
     { id: 'p4', name: 'Bukayo Saka', photo: '/path/to/saka.png', team: premierLeagueTeams.ars, position: 'FW', value: 12 },
  ],
   mostCleanSheets: [
     { id: 'p5', name: 'Alisson Becker', photo: '/path/to/alisson.png', team: premierLeagueTeams.liv, position: 'GK', value: 14 },
     { id: 'p6', name: 'Ederson', photo: '/path/to/ederson.png', team: premierLeagueTeams.mci, position: 'GK', value: 13 },
  ],
  table: [
    { position: 1, team: premierLeagueTeams.mci, played: 30, won: 25, drawn: 3, lost: 2, goalsFor: 80, goalsAgainst: 20, goalDifference: 60, points: 78 },
    { position: 2, team: premierLeagueTeams.liv, played: 30, won: 24, drawn: 4, lost: 2, goalsFor: 75, goalsAgainst: 25, goalDifference: 50, points: 76 },
    { position: 3, team: premierLeagueTeams.ars, played: 30, won: 23, drawn: 5, lost: 2, goalsFor: 70, goalsAgainst: 22, goalDifference: 48, points: 74 },
  ],
  history: [
    { season: '2023/24', winner: premierLeagueTeams.mci, topScorer: { id: 'p1', name: 'Erling Haaland', photo: '/public/Erling Haaland.png', team: premierLeagueTeams.mci, position: 'FW', value: 36 } },
  ],
};

const laLigaMockData: LeagueData = {
  name: 'La Liga',
  logo: '/public/LaLiga.svg',
  themeColor: '#ee8707', // La Liga orange
  upcomingMatches: [
     { id: 'll_m1', homeTeam: laLigaTeams.rma, awayTeam: laLigaTeams.bar, date: '2025-04-06', time: '21:00', venue: 'Bernabeu', status: 'scheduled' },
  ],
  topScorers: [
      { id: 'p7', name: 'Jude Bellingham', photo: '/public/Jude Bellingham.png', team: laLigaTeams.rma, position: 'MF', value: 18 },
  ],
  topAssists: [],
  mostCleanSheets: [],
  table: [
      { position: 1, team: laLigaTeams.rma, played: 30, won: 26, drawn: 3, lost: 1, goalsFor: 70, goalsAgainst: 15, goalDifference: 55, points: 81 },
      { position: 2, team: laLigaTeams.bar, played: 30, won: 24, drawn: 4, lost: 2, goalsFor: 65, goalsAgainst: 20, goalDifference: 45, points: 76 },
  ],
  history: [
      { season: '2023/24', winner: laLigaTeams.rma },
  ],
};

const serieAMockData: LeagueData = {
  name: 'Serie A',
  logo: '/public/Lega Serie A.svg',
  themeColor: '#024494', // Serie A blue
  upcomingMatches: [],
  topScorers: [],
  topAssists: [],
  mostCleanSheets: [],
  table: [
      { position: 1, team: serieATeams.int, played: 30, won: 27, drawn: 2, lost: 1, goalsFor: 85, goalsAgainst: 18, goalDifference: 67, points: 83 },
  ],
  history: [],
};

const bundesligaMockData: LeagueData = {
  name: 'Bundesliga',
  logo: '/public/Bundesliga.svg',
  themeColor: '#d20515', // Bundesliga red
  upcomingMatches: [],
  topScorers: [],
  topAssists: [],
  mostCleanSheets: [],
  table: [
      { position: 1, team: bundesligaTeams.bay, played: 30, won: 25, drawn: 3, lost: 2, goalsFor: 90, goalsAgainst: 25, goalDifference: 65, points: 78 },
  ],
  history: [],
};

const ligue1MockData: LeagueData = {
  name: 'Ligue 1',
  logo: '/public/Ligue 1 Uber Eats.svg',
  themeColor: '#dae025', // Ligue 1 yellow/gold
  upcomingMatches: [],
  topScorers: [
      { id: 'p8', name: 'Kylian Mbappe', photo: '/public/Kylian Mbappe.png', team: ligue1Teams.psg, position: 'FW', value: 28 },
  ],
  topAssists: [],
  mostCleanSheets: [],
  table: [
      { position: 1, team: ligue1Teams.psg, played: 30, won: 28, drawn: 1, lost: 1, goalsFor: 95, goalsAgainst: 15, goalDifference: 80, points: 85 },
  ],
  history: [],
};


// --- Export all league data ---
export const allLeaguesData: Record<string, LeagueData> = {
  'premier-league': premierLeagueMockData,
  'la-liga': laLigaMockData,
  'serie-a': serieAMockData,
  'bundesliga': bundesligaMockData,
  'ligue-1': ligue1MockData,
};

export type { LeagueData, Match, PlayerStat, LeagueTableRow, LeagueHistoryItem, Team, Player, LeagueMetadata }; // Export types

export const leagueMetadata: Record<string, LeagueMetadata> = {
  'PL': {
    id: 'PL',
    name: 'Premier League',
    code: 'PL',
    logo: '/Premier League.svg',
    themeColor: '#37003C'
  },
  'PD': {
    id: 'PD',
    name: 'La Liga',
    code: 'PD',
    logo: '/LaLiga.svg',
    themeColor: '#FF4B7D'
  },
  'SA': {
    id: 'SA',
    name: 'Serie A',
    code: 'SA',
    logo: '/Lega Serie A.svg',
    themeColor: '#024494'
  },
  'BL1': {
    id: 'BL1',
    name: 'Bundesliga',
    code: 'BL1',
    logo: '/Bundesliga.svg',
    themeColor: '#D20515'
  },
  'FL1': {
    id: 'FL1',
    name: 'Ligue 1',
    code: 'FL1',
    logo: '/Ligue 1 Uber Eats.svg',
    themeColor: '#DED531'
  },
  'CL': {
    id: 'CL',
    name: 'Champions League',
    code: 'CL',
    logo: '/UCL.svg',
    themeColor: '#1D3072'
  },
  'EL': {
    id: 'EL',
    name: 'Europa League',
    code: 'EL',
    logo: '/UEL.svg',
    themeColor: '#FF6900'
  }
};