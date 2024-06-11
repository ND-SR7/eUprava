import Warrant from "../../../models/Court/Warrant";
import WarrantCardStyled from "./WarrantCard/WarrantCard.styled";
import WarrantListStyled from "./WarrantList.styled";

interface WarrantListProps {
  warrants: Warrant[];
};
  
const WarrantList = ({warrants}: WarrantListProps) => {
  const content = warrants.map(warrant => 
    <WarrantCardStyled>
      <h1>{warrant.trafficViolation}</h1>
      <h6>{warrant.id}</h6>
      <h3>Issued on: {warrant.issuedOn.replace("T", " ").replace("Z", "")}</h3>
    </WarrantCardStyled>
  );

  return (
    <WarrantListStyled>{content}</WarrantListStyled>
  );
};

export default WarrantList;
