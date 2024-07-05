import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { CheckVehicleRegistration } from '../../../services/PoliceService';
import { FormContainer, Input, SubmitButton, StyledHeading } from './CheckVehicleRegistrationForm.styled';

interface CheckVehicleRegistrationFormProps {
  closeModal: () => void;
}

const CheckVehicleRegistrationForm: React.FC<CheckVehicleRegistrationFormProps> = ({ closeModal }) => {
  const [jmbg, setJmbg] = useState('');
  const [location, setLocation] = useState('');
  const [platesNumber, setPlatesNumber] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const data = {
      jmbg,
      location,
      platesNumber,
    };

    try {
      const response = await CheckVehicleRegistration(data);
      toast.success(response.message || "Vehicle registration checked successfully");
      navigate('/home/police');
      closeModal();
    } catch (error: any) {
      console.error(error);
      toast.error('Failed to check vehicle registration, please verify the details and try again.');
    }
  };

  return (
    <FormContainer>
      <StyledHeading>Check Vehicle Registration</StyledHeading>
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
        <Input
          type="text"
          placeholder="Plates Number"
          value={platesNumber}
          onChange={(e) => setPlatesNumber(e.target.value)}
          required
        />
        <SubmitButton buttonType="submit" label="Check Vehicle Registration" />
      </form>
    </FormContainer>
  );
};

export default CheckVehicleRegistrationForm;
