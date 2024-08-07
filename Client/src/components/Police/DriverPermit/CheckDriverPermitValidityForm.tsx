import React, { useState } from "react";
import toast from "react-hot-toast";
import HeadingStyled from "../../Shared/Heading/Heading.styled";
import Button from "../../Shared/Button/Button";
import { useNavigate } from "react-router-dom";
import { CheckDriverPermitValidation } from "../../../services/PoliceService";
import { FormContainer, Input } from "./CheckDriverPermitValidity.styled";

interface CheckDriverPermitValidityFormProps {
  closeModal: () => void;
}

const CheckDriverPermitValidityForm: React.FC<CheckDriverPermitValidityFormProps> = ({ closeModal }) => {
  const [jmbg, setJmbg] = useState("");
  const [location, setLocation] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    const data = {
      jmbg,
      location,
    };

    try {
      const response = await CheckDriverPermitValidation(data);
      toast.success(response.message || "Driver permit status checked successfully");
      navigate("/home/police");
      closeModal();
    } catch (error: any) {
      console.error(error);
      toast.error("Failed to check driver permit status, check driver JMBG");
    }
  };

  return (
    <FormContainer>
      <HeadingStyled>Check Driver Permit Validity</HeadingStyled>
      <form onSubmit={handleSubmit}>
        <Input
          type="text"
          placeholder="JMBG"
          value={jmbg}
          onChange={(e) => setJmbg(e.target.value)}
          required
        />
        <Input
          type="text"
          placeholder="Location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          required
        />
        <Button buttonType="submit" label="Check Driver Permit Validity" />
      </form>
    </FormContainer>
  );
};

export default CheckDriverPermitValidityForm;
