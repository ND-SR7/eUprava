import React, { useState } from "react";
import toast from "react-hot-toast";
import HeadingStyled from "../../Shared/Heading/Heading.styled";
import Button from "../../Shared/Button/Button";
import { useNavigate } from "react-router-dom";
import { checkVehicleTire } from "../../../services/PoliceService";
import Select from "../../Shared/Select/Select";
import { FormContainer, Input } from "./CheckVehicleTireForm.stled";

interface CheckVehicleTireFormProps {
    closeModal: () => void;
}

const CheckVehicleTireForm: React.FC<CheckVehicleTireFormProps> = ({ closeModal }) => {
    const [tireType, setTireType] = useState("WINTER");
    const [jmbg, setJmbg] = useState("");
    const [location, setLocation] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e: any) => {
        e.preventDefault();

        const data = {
            tireType,
            jmbg,
            location,
        };

        try {
            const response = await checkVehicleTire(data);
            console.log(response);
            toast.success(response.data.message || "Vehicle tire checked successfully");
            navigate("/home/police");
            closeModal();
        } catch (error: any) {
            console.error(error);
            toast.error("Failed to check vehicle tire, check driver JMBG");
        }
    };

    const handleTireTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setTireType(e.target.value);
    };

    return (
        <FormContainer>
            <HeadingStyled>Check Vehicle Tire</HeadingStyled>
            <form onSubmit={handleSubmit}>
                <Input
                    type="text"
                    placeholder="JMBG"
                    value={jmbg}
                    onChange={(e) => setJmbg(e.target.value)}
                    required
                />
                <Select
                    value={tireType}
                    onChange={handleTireTypeChange}
                    required={true}
                />
                <Input
                    type="text"
                    placeholder="Location"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                    required
                />
                <Button buttonType="submit" label="Check Vehicle Tire" />
            </form>
        </FormContainer>
    );
};

export default CheckVehicleTireForm;
