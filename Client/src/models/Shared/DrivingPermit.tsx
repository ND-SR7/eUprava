interface TrafficPermit {
    id: string; 
    number: string;
    issuedDate: string;
    expirationDate: string; 
    approved: boolean;
    person: string;
};

export default TrafficPermit;