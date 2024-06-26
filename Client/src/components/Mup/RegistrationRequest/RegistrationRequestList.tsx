import Registration from "../../../models/Shared/Registration";
import { approveRegistrationRequest, declineRegistrationRequest } from "../../../services/MupService";
import Button from "../../Shared/Button/Button";
import RegistrationRequestCardStyled from "./RegistratonRequestCard/RegistrationRequestCard.styled";

interface RegistrationRequestsProps {
    registrationRequests: Registration[];
    closeModal: () => void;
    refreshRequests: () => void;
}

const RegistrationRequestCard = ({ registrationRequests, closeModal, refreshRequests }: RegistrationRequestsProps) => {
    const content = registrationRequests.map(request => (
        <RegistrationRequestCardStyled key={request.registrationNumber}>
            <h1>Number: {request.registrationNumber}</h1>
            <h1>Issued date: {request.issuedDate.replace("T", " ").replace("Z", "")}</h1>
            <h1>Expiration date: {request.expirationDate.replace("T", " ").replace("Z", "")}</h1>
            <Button
                label="Approve"
                buttonType="button"
                onClick={() => approveRegistrationRequest(request.registrationNumber, request.vehicleID, request.owner, closeModal, refreshRequests)}
            />
            <Button
                label="Delete"
                buttonType="button"
                onClick={() => declineRegistrationRequest(request.registrationNumber, closeModal, refreshRequests)}
            />
        </RegistrationRequestCardStyled>
    ));

    return (
        <RegistrationRequestCardStyled>{content}</RegistrationRequestCardStyled>
    );
};

export default RegistrationRequestCard;
