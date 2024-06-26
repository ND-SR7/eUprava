import DrivingBansCardStyled from "./DrivingBansCard/DrivingBanCard.styled";
import DrivingBanListStyled from "./DrivingBanList.styled"; 
import DrivingBan from "../../../models/Mup/DrivingBans";

interface DrivingBansProps {
  drivingBans: DrivingBan[];
};

const DrivingBanList = ({drivingBans}: DrivingBansProps) => {
  const content = drivingBans.map(drivingBan =>
    <DrivingBansCardStyled>
      <h1>Driving ban reason: {drivingBan.reason}</h1>
      <h3>Duration: {drivingBan.duration.replace("T", " ").replace("Z", "")}</h3>
    </DrivingBansCardStyled>
  );

  return (
    <DrivingBanListStyled>{content}</DrivingBanListStyled>
  );
};

export default DrivingBanList;
