import CourtHearing from "../../../models/Court/CourtHearing";
import HearingCardStyled from "./HearingCard/HearingCard.styled";
import HearingListStyled from "./HearingList.styled";

interface HearingsProps {
  hearings: CourtHearing[];
};

const HearingList = ({hearings}: HearingsProps) => {
  const content = hearings.map(hearing =>
    <HearingCardStyled>
      <h1>{hearing.reason}</h1>
      <h6>{hearing.id}</h6>
      <h3>Date and time: {hearing.dateTime.replace("T", " ").replace("Z", "")}</h3>
      <p><b>Court: {hearing.court}</b></p>
    </HearingCardStyled>
  );

  return (
    <HearingListStyled>{content}</HearingListStyled>
  );
};

export default HearingList;
