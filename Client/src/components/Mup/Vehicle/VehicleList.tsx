import VehicleDTO from "../../../models/Mup/VehicleDetails";
import VehiclesCardStyled from "./VehiclesCard/VehicleCard.styled";
import VehiclesListStyled from "./VehicleList.styled";
import Button from "../../Shared/Button/Button";
import { handleRegister } from "../../../services/MupService";

interface VehiclesProps {
    vehicles: VehicleDTO[];
}

const VehicleList = ({ vehicles }: VehiclesProps) => {
    const content = vehicles.map(vehicle => (
        <VehiclesCardStyled key={vehicle.id}>
            <h1>Brand: {vehicle.brand}</h1>
            <h1>Model: {vehicle.model}</h1>
            <h1>Year: {vehicle.year}</h1>
            <h1>Registration: {vehicle.registration.approved === false ? "your vehicle is not registered" : vehicle.registration.registrationNumber}</h1>
            <h1>Plates: {vehicle.registration.approved === false ? "your vehicle is not registered" : vehicle.registration.plates}</h1>
            {vehicle.registration.approved === false && (
                <Button
                    label="Register Vehicle"
                    buttonType="button"
                    onClick={() => handleRegister(vehicle.id!)}
                />
            )}
        </VehiclesCardStyled>
    ));

    return (
        <VehiclesListStyled>{content}</VehiclesListStyled>
    );
};

export default VehicleList;
