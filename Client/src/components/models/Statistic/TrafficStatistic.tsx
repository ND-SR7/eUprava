import Vehicle from "../Shared/Vehicle";

type TrafficStatistic = {
    id: string;
    date: string;
    region: string;
    year: number;
    month: number;
    violationType: string;
    vehicles: Vehicle[];
};
  
export default TrafficStatistic;
  