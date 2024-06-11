import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingCourt } from "../services/PingService";
import toast from "react-hot-toast";
import Modal from "../components/Shared/Modal/Modal";
import { getHearingsByJMBG, getSuspensionByJMBG, getWarrantsByJMBG } from "../services/CourtService";
import CourtHearing from "../models/Court/CourtHearing";
import SuspensionModel from "../models/Court/Suspension";
import Warrant from "../models/Court/Warrant";
import HearingList from "../components/Court/Hearing/HearingList";
import Suspension from "../components/Court/Suspension/Suspension";
import WarrantList from "../components/Court/Warrant/WarrantList";

const CourtPage = () => {
  const navigate = useNavigate();
  const [modalVisible, setModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);
  const [modalHeading, setModalHeading] = useState("");
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
    pingCourt().then(() => {
      toast.success("Good connection with court service");
    }).catch((error) => {
      toast.error("No connection to court service");
      console.error(error);
    });
  };

  const checkHearings = () => {
    getHearingsByJMBG().then((hearings: CourtHearing[]) => {
      if (hearings.length > 0) {
        setModalContent(<HearingList hearings={hearings} />);
        setModalHeading("Your hearings");
      } else {
        setModalContent("");
        setModalHeading("No court hearings for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch hearings");
      console.error(error);
    });
  };

  const checkSuspensions = () => {
    getSuspensionByJMBG().then((suspension: SuspensionModel) => {
      if (suspension.id !== undefined) {
        setModalContent(<Suspension suspension={suspension} />);
        setModalHeading("Your suspension");
      } else {
        setModalContent("");
        setModalHeading("No suspension issued for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch suspension");
      console.error(error);
    });
  };

  const checkWarrants = () => {
    getWarrantsByJMBG().then((warrants: Warrant[]) => {
      if (warrants !== null && warrants.length > 0) {
        setModalContent(<WarrantList warrants={warrants} />);
        setModalHeading("Your warrants");
      } else {
        setModalContent("");
        setModalHeading("No warrants issued for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch warrants");
      console.error(error);
    });
  };

  const closeModal = () => {
    setModalVisible(false);
  };
  
  return (
    <>
      <HeadingStyled>Court</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <br />
      <Button buttonType="button" label="Check for hearings" onClick={() => checkHearings()} />
      <Button buttonType="button" label="Check for suspensions" onClick={() => checkSuspensions()} />
      <Button buttonType="button" label="Check for warrants" onClick={() => checkWarrants()} />
      <br />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />

      <Modal
        heading={modalHeading}
        content={modalContent}
        isVisible={modalVisible}
        onClose={closeModal}
      />
    </>
  );
};

export default CourtPage;
