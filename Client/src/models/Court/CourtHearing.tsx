type CourtHearing = {
  id: string;
  reason: string;
  dateTime: string;
  court: string;
  person?: string;
  legalEntity?: string;
};

export default CourtHearing;
