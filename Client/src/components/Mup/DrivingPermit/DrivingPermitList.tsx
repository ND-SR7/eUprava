import DrivingPermit from "../../../models/Shared/DrivingPermit";
import { handleRequestPermit } from "../../../services/MupService";
import Button from "../../Shared/Button/Button";
import DrivingPermitCardStyled from "./DrivingPermitCard/DrivingPermitCard";
import DrivingPermitStyled from "./DrivingPermitList.styled";

interface DrivingPermitProps {
    drivingPermits: DrivingPermit[];
};

const DrivingPermitCard = ({ drivingPermits }: DrivingPermitProps) => {
    const content = drivingPermits.map(drivingPermit => (
        <DrivingPermitCardStyled key={drivingPermit.number}>
            <h1>Number: {drivingPermit.number}</h1>
            <h1>Issued date: {drivingPermit.issuedDate.replace("T", " ").replace("Z", "")}</h1>
            <h1>Expiration date: {drivingPermit.expirationDate.replace("T", " ").replace("Z", "")}</h1>
        </DrivingPermitCardStyled>
    ));

    return (
        <DrivingPermitStyled>
                        <Button
                label="Request Driving Permit"
                buttonType="button"
                onClick={handleRequestPermit}
            />
            {content}
        </DrivingPermitStyled>
    );
};

export default DrivingPermitCard;