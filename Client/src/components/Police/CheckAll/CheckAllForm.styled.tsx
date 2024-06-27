import styled from 'styled-components';
import HeadingStyled from '../../Shared/Heading/Heading.styled';
import Button from '../../Shared/Button/Button';

export const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  width: 100%;
  max-width: 600px; /* Adjust the width as needed */
  margin: 0 auto; /* Center align */
  background-color: #ffffff;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
  border-radius: 8px;
`;

export const Input = styled.input`
  margin: 10px 0;
  padding: 12px;
  width: 100%;
  max-width: 300px;
  border: 1px solid #ddd;
  border-radius: 4px;
`;

export const Select = styled.select`
  margin: 10px 0;
  padding: 12px;
  width: 100%;
  max-width: 300px;
  border: 1px solid #ddd;
  border-radius: 4px;
`;

export const StyledHeading = styled(HeadingStyled)`
  font-size: 1.5em;
  color: #009879;
  margin-bottom: 20px;
`;

export const SubmitButton = styled(Button)`
  background-color: #009879;
  color: #ffffff;
  padding: 12px 24px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:hover {
    background-color: #007c6e;
  }
`;
