package updateserver

import (
	"fmt"

	"github.com/insys-icom/autoupdate/pkg/repository"
)

func hasActiveDownloadList(dao repository.DAO, tenant, target string) (bool, error) {
	update, err := dao.GetUpdateTargetBySlug(tenant, target)
	if err != nil {
		return false, err
	}

	if update.Enabled && update.DownloadListID.Valid {
		return true, nil
	}

	return false, nil
}

func createDownloadList(dao repository.DAO, prefix, tenant, target string) (string, error) {
	var result string

	// Get the addressed update target
	updateTarget, err := dao.GetUpdateTargetBySlug(tenant, target)
	if err != nil {
		return "", err
	}

	// Check if update target is active and contains a download list
	if !updateTarget.Enabled || !updateTarget.DownloadListID.Valid {
		return "", fmt.Errorf("update target has no active download list")
	}

	// Fetch the download list
	downloadList, _ := dao.GetDownloadListByID(tenant,
		updateTarget.DownloadListID.Int64)

	// Add all line items to download list output
	lineItems, _ := dao.FindDownloadListLineItemsByDownloadListID(tenant,
		downloadList.ID)
	for _, lineItem := range lineItems {
		result += fmt.Sprintf("%s;%s/%s/%s/packages/%d.tar\n",
			lineItem.SerialNumber, prefix, tenant, target,
			lineItem.UpdatePacketID)
	}

	// Add default (*) line item to download list output
	if downloadList.DefaultUpdatePacketID.Valid {
		result += fmt.Sprintf("*;%s/%s/%s/packages/%d.tar\n", prefix, tenant,
			target, downloadList.DefaultUpdatePacketID.Int64)
	}

	// Check if the download list exposes any valid download item,
	// if not return an error
	if len(lineItems) == 0 && !downloadList.DefaultUpdatePacketID.Valid {
		return "", fmt.Errorf("update target contains no valid download list")
	}

	return result, nil
}
