import SuspensionModel from "../../../models/Court/Suspension";
import SuspensionStyled from "./Suspension.styled";

interface SuspensionProps {
  suspension: SuspensionModel;
};
  
const Suspension = ({suspension}: SuspensionProps) => {
  return (
    <SuspensionStyled>
      <h6>{suspension.id}</h6>
      <h1>From: {suspension.from.replace("T", " ").replace("Z", "")}</h1>
      <h1>To: {suspension.to.replace("T", " ").replace("Z", "")}</h1>
    </SuspensionStyled>
  );
};

export default Suspension;
