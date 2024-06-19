import React, { useState } from "react";
import SelectStyled from "./SelectStyled.styled";

interface SelectProps {
  value: string;
  onChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
  required?: boolean;
}

const Select: React.FC<SelectProps> = ({ value, onChange, required }) => {
  return (
    <SelectStyled value={value} onChange={onChange} required={required}>
      <option value="WINTER">Winter</option>
      <option value="SUMMER">Summer</option>
    </SelectStyled>
  );
};

export default Select;
