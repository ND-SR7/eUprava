import DrivingPermit from "../../../models/Shared/DrivingPermit";
import { approveDrivingPermitRequest, declineDrivingPermitRequest } from "../../../services/MupService";
import Button from "../../Shared/Button/Button";
import DrivingPermitRequestCardStyled from "./DrivingPermitRequestCard/DrivingPermitRequestCard.styled";

interface DrivingPermitRequestsProps {
    drivingPermitRequests: DrivingPermit[];
    closeModal: () => void;
    refreshRequests: () => void; 
}

const DrivingPermitsRequestCard = ({ drivingPermitRequests, closeModal, refreshRequests }: DrivingPermitRequestsProps) => {
    const content = drivingPermitRequests.map(request => (
        <DrivingPermitRequestCardStyled key={request.id}>
            <h1>Number: {request.number}</h1>
            <h1>Issued date: {request.issuedDate.replace("T", " ").replace("Z", "")}</h1>
            <h1>Expiration date: {request.expirationDate.replace("T", " ").replace("Z", "")}</h1>
            <Button
                label="Approve"
                buttonType="button"
                onClick={() => approveDrivingPermitRequest(request.id, request.number, request.issuedDate, request.expirationDate, request.person, closeModal, refreshRequests)}
            />
            <Button
                label="Decline"
                buttonType="button"
                onClick={() => declineDrivingPermitRequest(request.id, closeModal, refreshRequests)}
            />
        </DrivingPermitRequestCardStyled>
    ));

    return (
        <DrivingPermitRequestCardStyled>{content}</DrivingPermitRequestCardStyled>
    );
};

export default DrivingPermitsRequestCard;
