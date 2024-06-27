type CourtHearing = {
  id?: string;
  reason: string;
  dateTime: string;
  court: string;
  person?: string;
  legalEntity?: string;
};

export type RescheduleCourtHearing = {
	hearingID: string;
	dateTime: string;
};

export default CourtHearing;
