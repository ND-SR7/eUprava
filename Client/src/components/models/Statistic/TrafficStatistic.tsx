type Vehicle = {
    id: string;
    brand: string;
    model: string;
    year: number;
    registration: string;
    plates: string;
    owner: string;
};
  
type TrafficStatistics = {
    id: string;
    date: string;
    region: string;
    year: number;
    month: number;
    violationType: string;
    vehicles: Vehicle[];
};
  
export type { TrafficStatistics };
  