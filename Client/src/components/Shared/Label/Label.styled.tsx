import styled from "styled-components";

const LabelStyled = styled.label`
  z-index: 2;
  color: ${(props) => props.theme.colors.accent};
  font-weight: ${(props) => props.theme.fontWeights.bold};
`;

export default LabelStyled;
