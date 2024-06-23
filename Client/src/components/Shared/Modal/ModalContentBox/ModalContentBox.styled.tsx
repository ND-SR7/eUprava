import styled from "styled-components";

const ModalContentBoxStyled = styled.div`
  background-color: ${(props) => props.theme.colors.background};
  margin: auto;
  padding: ${(props) => props.theme.paddings.large};
  border: 1px solid #888;
  width: 80%;
  max-height: 80vh; /* Maksimalna visina */
  overflow-y: auto; /* Omogućava vertikalno skrolovanje */
  position: relative;
`;

export default ModalContentBoxStyled;
