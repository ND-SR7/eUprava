interface RegistrationDetails {
  registrationNumber: string;
  issuedDate: string; 
  expirationDate: string; 
  vehicleID: string;
  owner: string;
  plates: string;
  approved: boolean;
};

interface PlatesDetails {
  registrationNumber: string;
  platesNumber: string;
  plateType: string;
  owner: string;
  vehicleID: string;
};

interface VehicleDTO {
  id?: string; 
  brand: string;
  model: string;
  year: number;
  registration: RegistrationDetails;
  plates: PlatesDetails;
  owner: string;
};

export default VehicleDTO;